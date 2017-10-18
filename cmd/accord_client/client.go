package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
	"golang.org/x/sys/unix"

	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord"
	"github.com/mistsys/accord/db"
	"github.com/mistsys/accord/id"
	"github.com/mistsys/accord/protocol"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultAddress = "localhost:50051"
	defaultName    = "NewHost"
	nanosInSecond  = 1000000000
	defaultSalt    = "hUYh5x4N2DOnTIce"
)

type stringSlice []string

// Implement the Value interface
func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Value() []string {
	return *s
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// calculate the latency from the metadata
// int64 in case the servers are out of sync
func Latency(metadata *protocol.ReplyMetadata, curTime time.Time) (time.Duration, time.Duration) {
	respTimeNsecs := int64(metadata.ResponseTime.Seconds)*nanosInSecond + int64(metadata.ResponseTime.Nanos)
	reqTimeNsecs := int64(metadata.RequestTime.Seconds)*nanosInSecond + int64(metadata.RequestTime.Nanos)
	serverNsecs := respTimeNsecs - reqTimeNsecs
	totalDuration := curTime.Sub(time.Unix(0, reqTimeNsecs))
	return time.Duration(serverNsecs), totalDuration
}

type Host struct {
	c            protocol.CertClient
	dryrun       bool
	pskStore     accord.PSKStore
	salt         string
	deploymentId string
	keysDir      string
	hostnames    []string
	uuid         []byte
}

// For now the response doesn't do a proper challenge auth
func (h *Host) Authenticate(ctx context.Context) (string, error) {
	keyId, err := id.KeyID(h.deploymentId, h.salt)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get the KeyId based on deploymentId")
	}
	//log.Printf("keyId: %d", keyId)
	aesgcm := accord.InitAESGCM(h.pskStore)
	// TODO: validate the instance doc and send that to the server to keep in records
	encrypted, err := aesgcm.Encrypt([]byte(`THIS IS WHERE INSTANCE DOC WOULD GO`), keyId)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to encrypt the message")
	}

	req := &protocol.HostAuthRequest{
		RequestTime: ptypes.TimestampNow(),
		AuthInfo:    encrypted,
	}

	resp, err := h.c.HostAuth(ctx, req)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to send the authentication challenge")
	}
	h.uuid = resp.AuthResponse
	return string(h.uuid), nil
}

// It's always in the future, so not giving the user start time option
func (h *Host) RequestCerts(ctx context.Context, duration time.Duration) error {
	if h.keysDir == "" {
		return errors.New("keysDir isn't set, don't know where to read the public keys from")
	}
	if unix.Access(h.keysDir, unix.W_OK) != nil {
		return fmt.Errorf("Directory %s isn't writable, aborting before requesting certs", h.keysDir)
	}
	files, err := listPubKeysInDir(h.keysDir)
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
			Id:          h.uuid,
			Hostnames:   h.hostnames,
		}
		resp, err := h.c.HostCert(ctx, certRequest)
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

func certPath(pubKeyPath string) string {
	base := filepath.Base(pubKeyPath)
	dir := filepath.Dir(pubKeyPath)
	prefix := strings.Split(base, ".")[0]
	return path.Join(dir, prefix+"-cert.pub")
}

type User struct {
	c              protocol.CertClient
	keysDir        string
	username       string
	remoteUsername string
	principals     []string
	token          *oauth2.Token
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

// Returns full path for all the public keys in the directory given
func listPubKeysInDir(dir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to enumerate files from %s", dir)
	}
	files := []string{}
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), ".pub") && !strings.Contains(fileInfo.Name(), "cert") {
			files = append(files, path.Join(dir, fileInfo.Name()))
		}
	}
	return files, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	address := flag.String("server", defaultAddress, "The grpc server to contact")
	task := flag.String("task", "hostcert", "Task to run, whether you want host cert, user cert")
	psk := flag.String("psk", "", "Pre-shared key to use, this will be removed at some point in future")
	insecure := flag.Bool("insecure", false, "Is this for development, and disable TLS?")
	dryrun := flag.Bool("dryrun", false, "Is this for testing?")
	nowebserver := flag.Bool("nowebserver", false, "Start webserver to ge the prompt or not")
	//overwrite := flag.Bool("o", false, "Overwrite the files")
	//duration := flag.Duration("duration", 1*time.Hour, "How long to get the cert for")
	deploymentId := flag.String("deploymentId", "", "ID to use for authenticating with the server")
	hostKeysPath := flag.String("hostkeys", "/etc/ssh", "Where to read the public keys")
	remoteUsername := flag.String("remoteusername", "", "What remote username to allow")
	serverCert := flag.String("cert", "", "Server cert to use")
	hostSalt := flag.String("hostsalt", defaultSalt, "Randomly generated string to prefix requests when creating host requests")
	userKeysPath := flag.String("userkeys", "", "Where to find the user's public keys")
	domain := flag.String("domain", "mistsys.com", "Google Apps Domain to validate for")
	var (
		hostnames  = stringSlice{}
		principals = stringSlice{}
	)
	flag.Var(&hostnames, "host", "Hostnames to sign for")
	flag.Var(&principals, "p", "Principals to validate for, these are usernames and server class, etc")
	flag.Parse()

	var (
		conn *grpc.ClientConn
		err  error
	)
	var done = make(chan struct{})
	// Set up a connection to the server.
	if *insecure {
		conn, err = grpc.Dial(*address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Unable to connect to %s: %v", *address, err)
		}
		defer conn.Close()
	} else {
		var creds credentials.TransportCredentials
		if *serverCert != "" {
			creds, err = credentials.NewClientTLSFromFile(*serverCert, "localhost")
			if err != nil {
				log.Fatalf("could not load tls cert: %s", err)
			}
			conn, err = grpc.Dial(*address, grpc.WithTransportCredentials(creds))
			if err != nil {
				log.Fatalf("Unable to connect to %s: %v", *address, err)
			}
		} else {
			// Use the operating system default root certificates.
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			}
			conn, err = grpc.Dial(*address, opts...)
			if err != nil {
				log.Fatalf("Unable to connect to %s: %v", *address, err)
			}
			//creds = credentials.NewClientTLSFromCert(nil, "")
		}

		defer conn.Close()
	}

	switch *task {
	case "hostcert":
		pskStore := db.NewSinglePSKStore(*deploymentId, *hostSalt, *psk)
		if pskStore == nil {
			log.Fatalf("cannot initialize the Simple PSK Store")
		}
		c := protocol.NewCertClient(conn)
		log.Println("Starting authentication for host")
		host := &Host{
			c:            c,
			dryrun:       *dryrun,
			pskStore:     pskStore,
			salt:         *hostSalt,
			deploymentId: *deploymentId,
			keysDir:      *hostKeysPath,
			hostnames:    hostnames,
		}

		uuid, err := host.Authenticate(context.Background())
		if err != nil {
			log.Fatalf("Failed to authenticate the host with cert server %s", err)
		}
		log.Printf("uuid: %s\n", uuid)
		err = host.RequestCerts(context.Background(), 30*24*time.Hour)
		if err != nil {
			log.Fatalf("Failed to get the certs %s", err)
		}
		close(done)
	case "usercert":
		c := protocol.NewCertClient(conn)
		log.Println("Starting authentication for user")
		googleAuth := &accord.GoogleAuth{
			Domain:        *domain,
			UseWebServer:  !*nowebserver,
			WebServerPort: 8091,
		}

		tok, err := googleAuth.Authenticate()
		if err != nil {
			log.Fatalf("Failed to authenticate user: %s", err)
		}

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		var keysDir = filepath.Join(usr.HomeDir, ".ssh")
		if *userKeysPath != "" {
			keysDir = *userKeysPath
		}

		if *remoteUsername == "" {
			log.Printf("remoteusername is empty, picking the same as current username: %s", usr.Name)
			remoteUsername = &usr.Name
		}

		//log.Printf("token: %v", tok)
		user := &User{
			c:              c,
			username:       usr.Username,
			remoteUsername: *remoteUsername,
			keysDir:        keysDir,
			token:          tok,
			principals:     principals.Value(),
		}
		ok, email, err := user.CheckAuthorization(context.Background())
		if err != nil {
			log.Fatalf("Failed to check authorization for the user: %s", err)
		}
		if !ok {
			log.Fatalf("Invalid state reached for user cert, cannot continue further")
		}

		err = user.RequestCerts(context.Background(), email, 24*time.Hour)
		if err != nil {
			log.Fatalf("Failed to get the certs %s", err)
		}
		close(done)
		//log.Fatalf("Not done yet")
	default:
		log.Fatalf("Don't know the task %s", *task)
	}
	<-done
	//log.Printf("Ping Response: %s. Latency: %s", r.Message, Latency(r.Metadata))
}
