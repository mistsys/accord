package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord"
	"github.com/mistsys/accord/cloud_metadata"
	"github.com/mistsys/accord/id"
	"github.com/mistsys/accord/protocol"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
)

// Host has all fields public because most of them are required to be set by
// the caller in one way or another anyway and duplicating the fields in
// another input object is just waste to only be used once
type Host struct {
	Client       protocol.CertClient
	Dryrun       bool
	PSKStore     accord.PSKStore
	Salt         string
	DeploymentId string
	KeysDir      string
	Hostnames    []string
	UUID         []byte
}

func NewHost(client protocol.CertClient) *Host {
	return &Host{
		Client: client,
	}
}

// For now the response doesn't do a proper challenge auth
func (h *Host) Authenticate(ctx context.Context) (string, error) {
	keyId, err := id.KeyID(h.DeploymentId, h.Salt)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get the KeyId based on deploymentId")
	}
	//log.Printf("keyId: %d", keyId)
	aesgcm := accord.InitAESGCM(h.PSKStore)

	cloud, err := cloud_metadata.CloudService()
	if err != nil {
		log.Printf("Failed to read metadata %s", err)
	}
	metadata := []byte("Unknown: test code")
	if cloud == cloud_metadata.AWS {
		instanceInfo, err := cloud_metadata.GetAWSInstanceInfo()
		if err != nil {
			return "", errors.Wrapf(err, "Failed to read AWS instance info")
		}
		metadata, err = json.Marshal(instanceInfo)
		if err != nil {
			return "", errors.Wrapf(err, "Failed to serialize the metadata")
		}
	} else {
		return "", errors.Wrapf(err, "Cloud %s not supported yet", cloud)
	}
	// sends the data to the server to keep for records
	encrypted, err := aesgcm.Encrypt(metadata, keyId)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to encrypt the message")
	}

	req := &protocol.HostAuthRequest{
		RequestTime: ptypes.TimestampNow(),
		AuthInfo:    encrypted,
	}

	resp, err := h.Client.HostAuth(ctx, req)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to send the authentication challenge")
	}
	h.UUID = resp.AuthResponse
	return string(h.UUID), nil
}

// It's always in the future, so not giving the user start time option
func (h *Host) RequestCerts(ctx context.Context, duration time.Duration) error {
	if h.KeysDir == "" {
		return errors.New("keysDir isn't set, don't know where to read the public keys from")
	}
	if unix.Access(h.KeysDir, unix.W_OK) != nil {
		return fmt.Errorf("Directory %s isn't writable, aborting before requesting certs", h.KeysDir)
	}
	files, err := listPubKeysInDir(h.KeysDir)
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

		certRequest := &protocol.HostCertRequest{
			RequestTime: ptypes.TimestampNow(),
			PublicKey:   contents,
			ValidFrom:   protoValidFrom,
			ValidUntil:  protoValidUntil,
			Id:          h.UUID,
			Hostnames:   h.Hostnames,
		}
		resp, err := h.Client.HostCert(ctx, certRequest)
		if err != nil {
			return errors.Wrapf(err, "Error when trying to get cert for %s", f)
		}
		certFileName := certPath(f)
		log.Printf("Writing to %s", certFileName)
		err = ioutil.WriteFile(certFileName, resp.HostCert, 0644)
		if err != nil {
			return errors.Wrapf(err, "Failed to write %s", certFileName)
		}
	}
	return nil
}

func (h *Host) UpdateUserCertAuthority(filePath string) error {
	resp, err := h.Client.PublicTrustedCA(context.Background(), &protocol.PublicTrustedCARequest{
		RequestTime: ptypes.TimestampNow(),
	})
	if err != nil {
		log.Fatalf("Failed to get the certs %s", err)
	}
	userCerts := [][]byte{}
	//fmt.Println("=== Host CAs ===")
	for _, userCA := range resp.UserCAs {
		userCerts = append(userCerts, userCA.PublicKey)
	}
	return updateUsersCertAuthority(filePath, userCerts)
}
