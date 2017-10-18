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
type CertAccorder struct {
	pskStore    accord.PSKStore
	certManager *accord.CertManager
	aesgcm      *accord.AESGCM
	domain      string
	authz       accord.Authz
}

func NewCertAccorder(pskStore accord.PSKStore, certManager *accord.CertManager, domain string, authz accord.Authz) *CertAccorder {
	return &CertAccorder{
		pskStore:    pskStore,
		certManager: certManager,
		aesgcm:      accord.InitAESGCM(pskStore),
		domain:      domain,
		authz:       authz,
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

func (s *CertAccorder) HostAuth(ctx context.Context, authRequest *protocol.HostAuthRequest) (*protocol.HostAuthResponse, error) {
	log.Println("Received host auth request")

	decrypted, err := s.aesgcm.Decrypt(authRequest.AuthInfo)
	if err != nil {
		// maybe wait until the deadline in Context and respond?
		// to handle for timing based attacks
		return nil, errors.Wrapf(err, "Failed to decrypt message.")
	}
	log.Printf("Decrypted message from host %s", string(decrypted))

	return &protocol.HostAuthResponse{
		Metadata:     replyMetadata(authRequest.GetRequestTime()),
		AuthResponse: makeUUID(),
	}, nil
}

func (s *CertAccorder) HostCert(ctx context.Context, certRequest *protocol.HostCertRequest) (*protocol.HostCertResponse, error) {
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

func (s *CertAccorder) UserAuth(ctx context.Context, userAuthRequest *protocol.UserAuthRequest) (*protocol.UserAuthResponse, error) {
	oauthToken, err := accord.OAuth2Token(userAuthRequest.Token)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot convert pb token to *oauth2.Token")
	}
	log.Printf("Received oauthToken: %#v", oauthToken)
	googleAuth := &accord.GoogleAuth{
		Domain: s.domain,
		Token:  oauthToken,
	}
	valid, email, err := googleAuth.ValidateToken(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to validate token")
	}
	return &protocol.UserAuthResponse{
		UserId: email,
		Valid:  valid,
	}, nil
}

func (s *CertAccorder) UserCert(ctx context.Context, certRequest *protocol.UserCertRequest) (*protocol.UserCertResponse, error) {

	validFrom, _ := ptypes.Timestamp(certRequest.ValidFrom)
	validUntil, _ := ptypes.Timestamp(certRequest.ValidUntil)

	authorizedPrincipals, err := s.authz.Authorized(certRequest.UserId, certRequest.AuthorizedPrincipals)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed authorization")
	}

	srq := &accord.CertSignRequest{
		PubKey:     certRequest.PublicKey,
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Id:         certRequest.RemoteUsername,
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

func (s *CertAccorder) PublicTrustedCA(context.Context, *protocol.PublicTrustedCARequest) (*protocol.PublicTrustedCAResponse, error) {
	panic("not implemented")
}

func (s *CertAccorder) Ping(ctx context.Context, req *protocol.PingRequest) (*protocol.PingResponse, error) {
	return &protocol.PingResponse{
		Metadata: replyMetadata(req.GetRequestTime()),
		Message:  "Hello " + req.GetName(),
	}, nil
}
