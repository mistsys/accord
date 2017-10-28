package certserver

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mistsys/accord"
	"github.com/mistsys/accord/protocol"
	"github.com/pkg/errors"
)

// I ran out of names to give
type AccordServer struct {
	pskStore       accord.PSKStore
	certManager    *accord.CertManager
	aesgcm         *accord.AESGCM
	googleClientId string
	domain         string
	authz          accord.Authz
}

func NewAccordServer(pskStore accord.PSKStore, certManager *accord.CertManager,
	googleClientId string,
	domain string, authz accord.Authz) *AccordServer {
	return &AccordServer{
		pskStore:       pskStore,
		certManager:    certManager,
		googleClientId: googleClientId,
		aesgcm:         accord.InitAESGCM(pskStore),
		domain:         domain,
		authz:          authz,
	}
}

func replyMetadata(reqTime *google_protobuf.Timestamp) *protocol.ReplyMetadata {
	return &protocol.ReplyMetadata{
		RequestTime:  reqTime,
		ResponseTime: ptypes.TimestampNow(),
	}
}

// it doesn't have to be spec compliant, just random enough to not collide
// we're talking about couple of hundred to thousand nodes
func makeUUID() []byte {
	b := make([]byte, 16)
	rand.Read(b)
	return []byte(fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
}

func (s *AccordServer) HostAuth(ctx context.Context, authRequest *protocol.HostAuthRequest) (*protocol.HostAuthResponse, error) {
	log.Println("Received host auth request")

	decrypted, nonce, sender, err := s.aesgcm.Decrypt(authRequest.AuthInfo)
	if err != nil {
		// maybe wait until the deadline in Context and respond?
		// to handle for timing based attacks
		return nil, errors.Wrapf(err, "Failed to decrypt message.")
	}
	log.Printf("Decrypted message from host %s", string(decrypted))

	uuid := makeUUID()
	// we have already established the server connection, no need for a new ID
	// additionally we can have a well known server key known by every client
	// but I don't see a lot of gain there
	// this will be going over an already-encrypted connection too
	encrypted, err := s.aesgcm.EncryptWithNonce(uuid, nonce, sender)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to encrypt with the nonce from client")
	}
	return &protocol.HostAuthResponse{
		Metadata:     replyMetadata(authRequest.GetRequestTime()),
		AuthResponse: encrypted,
	}, nil
}

func (s *AccordServer) HostCert(ctx context.Context, certRequest *protocol.HostCertRequest) (*protocol.HostCertResponse, error) {
	// TODO: validate the host ID properly
	validFrom, _ := ptypes.Timestamp(certRequest.ValidFrom)
	validUntil, _ := ptypes.Timestamp(certRequest.ValidUntil)
	srq := &accord.CertSignRequest{
		PubKey:     certRequest.PublicKey,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Id:         string(certRequest.Id),
		Serial:     1,
		Principals: certRequest.Hostnames,
	}
	hostCert, err := s.certManager.SignHostCert(srq)
	if err != nil {
		return &protocol.HostCertResponse{
			Metadata: replyMetadata(certRequest.GetRequestTime()),
		}, errors.Wrapf(err, "Failed to sign host cert for hostnames: %s", certRequest.Hostnames)
	}
	return &protocol.HostCertResponse{
		Metadata: replyMetadata(certRequest.GetRequestTime()),
		HostCert: hostCert,
	}, nil
}

func (s *AccordServer) UserAuth(ctx context.Context, userAuthRequest *protocol.UserAuthRequest) (*protocol.UserAuthResponse, error) {
	oauthToken, err := accord.OAuth2Token(userAuthRequest.Token)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot convert pb token to *oauth2.Token")
	}
	log.Printf("Received authentication token for user: %s", userAuthRequest.GetUsername())
	googleAuth := &accord.GoogleAuth{
		Domain:   s.domain,
		ClientId: s.googleClientId,
		Token:    oauthToken,
	}
	valid, email, err := googleAuth.ValidateToken(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to validate token")
	}
	log.Printf("Valid user: %s email: %s", userAuthRequest.GetUsername(), email)
	return &protocol.UserAuthResponse{
		UserId: email,
		Valid:  valid,
	}, nil
}

func (s *AccordServer) UserCert(ctx context.Context, certRequest *protocol.UserCertRequest) (*protocol.UserCertResponse, error) {

	validFrom, _ := ptypes.Timestamp(certRequest.ValidFrom)
	validUntil, _ := ptypes.Timestamp(certRequest.ValidUntil)

	authorizedPrincipals, err := s.authz.Authorized(certRequest.UserId, certRequest.AuthorizedPrincipals)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed authorization")
	}

	// TODO: increase serial every time user requests for a certificate
	// this requires some way to keep track of state, so shelving until
	// the project gets a database of some sort
	srq := &accord.CertSignRequest{
		PubKey:     certRequest.PublicKey,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Id:         certRequest.Username,
		Serial:     1,
		Principals: authorizedPrincipals,
	}

	userCert, err := s.certManager.SignUserCert(srq)
	if err != nil {
		return &protocol.UserCertResponse{
			Metadata: replyMetadata(certRequest.GetRequestTime()),
		}, errors.Wrapf(err, "Failed to sign user cert for %s", certRequest.Username)
	}
	return &protocol.UserCertResponse{
		Metadata: replyMetadata(certRequest.GetRequestTime()),
		UserCert: userCert,
	}, nil
}

// This doesn't populate the RevokedCerts yet
// TODO: try to use same data structure
func (s *AccordServer) PublicTrustedCA(ctx context.Context, trustedCARequest *protocol.PublicTrustedCARequest) (*protocol.PublicTrustedCAResponse, error) {
	hostCAs := s.certManager.HostCAs()
	userCAs := s.certManager.UserCAs()

	pbHostCAs := []*protocol.HostCA{}
	pbUserCAs := []*protocol.UserCA{}

	for _, h := range hostCAs {
		pbHostCAs = append(pbHostCAs, accord.ToHostCA(h))
	}

	for _, u := range userCAs {
		pbUserCAs = append(pbUserCAs, accord.ToUserCA(u))
	}

	return &protocol.PublicTrustedCAResponse{
		Metadata: replyMetadata(trustedCARequest.GetRequestTime()),
		HostCAs:  pbHostCAs,
		UserCAs:  pbUserCAs,
	}, nil

}

func (s *AccordServer) Ping(ctx context.Context, req *protocol.PingRequest) (*protocol.PingResponse, error) {
	return &protocol.PingResponse{
		Metadata: replyMetadata(req.GetRequestTime()),
		Message:  "Hello " + req.GetName(),
	}, nil
}
