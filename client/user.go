package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord"
	"github.com/mistsys/accord/protocol"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
	"golang.org/x/sys/unix"
)

type User struct {
	c              protocol.CertClient
	keysDir        string
	username       string
	remoteUsername string
	principals     []string
	token          *oauth2.Token
}

func NewUser(client protocol.CertClient) *User {
	return &User{
		c: client,
	}
}

func NewUserWithToken(client protocol.CertClient, token *oauth2.Token) *User {
	return &User{
		c:     client,
		token: token,
	}
}

func (u *User) SetUsername(username string) {
	u.username = username
}

func (u *User) SetRemoteUsername(username string) {
	u.remoteUsername = username
}

func (u *User) SetKeysDir(keysDir string) {
	u.keysDir = keysDir
}

func (u *User) SetPrincipals(principals []string) {
	u.principals = principals
}

// Validate the Oauth2.0 token against the server
// it needs to have been generated by the same clientID
// and have valid expiry date, etc
// Maybe use a nonce or something, but I can't think of
// reasonable things to send here, so boolean so that
// the backend doesn't get bad requests
func (u *User) CheckAuthorization(ctx context.Context) (bool, string, error) {
	if u.token == nil {
		return false, "", errors.New("Token isn't set, authenticate first with an OAuth2.0 service")
	}
	pbToken, err := accord.OAuthTokenPb(u.token)
	if err != nil {
		return false, "", errors.Wrapf(err, "Failed to convert *oauth2.Token to protobuf equivalent")
	}
	authRequest := &protocol.UserAuthRequest{
		RequestTime: ptypes.TimestampNow(),
		Username:    u.username,
		Token:       pbToken,
	}
	resp, err := u.c.UserAuth(ctx, authRequest)
	if err != nil {
		return false, "", errors.Wrapf(err, "Failed to authenticate user with server")
	}
	return resp.GetValid(), resp.GetUserId(), nil
}

func (u *User) RequestCerts(ctx context.Context, userId string, duration time.Duration) error {
	if len(u.principals) == 0 {
		return errors.New("No principals provided to request certificates for")
	}

	if u.keysDir == "" {
		return errors.New("keysDir isn't set, don't know where to read the public keys from")
	}
	if unix.Access(u.keysDir, unix.W_OK) != nil {
		return fmt.Errorf("Directory %s isn't writable, aborting before requesting certs", u.keysDir)
	}
	files, err := listPubKeysInDir(u.keysDir)
	if err != nil {
		return errors.Wrapf(err, "Failed list keys in host")
	}
	log.Printf("Found %d public keys that need to be signed", len(files))
	// 1. so that all the certs have same start time
	// 2. the server doesn't reject this for being too far in past
	// taking the number from Oauth2's implementation
	validFrom := time.Now().Add(10 * time.Second)
	validUntil := validFrom.Add(duration)
	for _, f := range files {
		contents, err := ioutil.ReadFile(f)
		if err != nil {
			return errors.Wrapf(err, "Failed to read %s", f)
		}
		_, _, _, _, err = ssh.ParseAuthorizedKey(contents)
		if err != nil {
			return errors.Wrapf(err, "%s doesn't look like a public key file", f)
		}
		protoValidFrom, err := ptypes.TimestampProto(validFrom)
		if err != nil {
			return errors.Wrapf(err, "can't make protobuf Timestamp for validFrom")
		}
		protoValidUntil, err := ptypes.TimestampProto(validUntil)
		if err != nil {
			return errors.Wrapf(err, "can't make protobuf Timestamp for validUntil")
		}

		var currentCert []byte
		certFileName := certPath(f)
		if _, err := os.Stat(certFileName); err != nil {
			// we should just create the file if it doesn't exist
			if !os.IsNotExist(err) {
				return errors.Wrapf(err, "Unknown error reading cert file %s")
			}
		} else {
			currentCert, err = ioutil.ReadFile(certFileName)
			if err != nil {
				return errors.Wrapf(err, "Found cert file %s but can't read")
			}
		}

		certRequest := &protocol.UserCertRequest{
			RequestTime:          ptypes.TimestampNow(),
			UserId:               userId,
			Username:             u.username,
			RemoteUsername:       u.remoteUsername,
			CurrentUserCert:      currentCert,
			PublicKey:            contents,
			ValidFrom:            protoValidFrom,
			ValidUntil:           protoValidUntil,
			AuthorizedPrincipals: u.principals,
		}
		resp, err := u.c.UserCert(ctx, certRequest)
		if err != nil {
			return errors.Wrapf(err, "Error when trying to get cert for %s", f)
		}

		log.Printf("Writing to %s", certFileName)
		err = ioutil.WriteFile(certFileName, resp.UserCert, 0644)
		if err != nil {
			return errors.Wrapf(err, "Failed to write %s", certFileName)
		}
	}
	return nil
}

func (u *User) UpdateHostCertAuthority(knownHostsFile string) error {
	resp, err := u.c.PublicTrustedCA(context.Background(), &protocol.PublicTrustedCARequest{
		RequestTime: ptypes.TimestampNow(),
	})
	if err != nil {
		log.Fatalf("Failed to get the certs %s", err)
	}
	hostCerts := [][]byte{}
	//fmt.Println("=== Host CAs ===")
	for _, hostCA := range resp.HostCAs {
		hostCerts = append(hostCerts, hostCA.PublicKey)
	}
	return updateKnownHostsCertAuthority(knownHostsFile, hostCerts)
}