package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mistsys/accord"
	"github.com/mistsys/accord/client"
	"github.com/mistsys/accord/db"
	"github.com/mistsys/accord/protocol"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultName = "NewHost"
	defaultSalt = "hUYh5x4N2DOnTIce"
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	address := flag.String("server", accord.DefaultServer, "The grpc server to contact")
	task := flag.String("task", "", "Task to run, whether you want host cert, user cert")
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
	googleClientId := flag.String("google.clientid", "", "Which Google Apps ClientID to use")
	googleClientSecret := flag.String("google.clientsecret", "", "Which Google Apps Client Secret to use, if not baked in already")
	domain := flag.String("domain", "mistsys.com", "Google Apps Domain to validate for")
	knownHostsFile := flag.String("knownhosts", "", "Known Hosts file, defaults to ~/.ssh/known_hosts")
	userCACertsFile := flag.String("userca", "", "Where the userca file should be, defaults to /etc/ssh/users_ca.pub")
	webserverPort := flag.Int("webserver.port", 8091, "Which port to run the auth webserver on")
	sshdFile := flag.String("sshdconfig", "", "SSHD Configuration file, defaults to /etc/ssh/sshd_config")
	var (
		hostnames  = stringSlice{}
		principals = stringSlice{}
	)
	flag.Var(&hostnames, "host", "Hostnames to sign for")
	flag.Var(&principals, "p", "Principals to validate for, these are usernames and server class, etc")
	flag.Parse()

	argv := flag.Args()
	fmt.Println(argv)

	// if no -task was given but there is one unparsed arg, it is the task
	if *task == "" && len(argv) >= 1 {
		task = &argv[0]
		argv = argv[1:]
	}

	var (
		conn *grpc.ClientConn
		err  error
	)
	// if the address was added by the build step, and has quotes
	serverAddress := accord.Unquote(*address)
	// this is very lazy way to keep the service running until job is done, in this kind of application
	// once I break it out further, it will change to something more sophisticated
	var done = make(chan struct{})
	// Set up a connection to the server.
	if *insecure {
		conn, err = grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Unable to connect to %s: %v", serverAddress, err)
		}
		defer conn.Close()
	} else {
		var creds credentials.TransportCredentials
		if *serverCert != "" {
			creds, err = credentials.NewClientTLSFromFile(*serverCert, "localhost")
			if err != nil {
				log.Fatalf("could not load tls cert: %s", err)
			}
			conn, err = grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))
			if err != nil {
				log.Fatalf("Unable to connect to %s: %v", serverAddress, err)
			}
		} else {
			// Use the operating system default root certificates.
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
			}
			conn, err = grpc.Dial(serverAddress, opts...)
			if err != nil {
				log.Fatalf("Unable to connect to %s: %v", serverAddress, err)
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
		host := &client.Host{
			Client:       c,
			Dryrun:       *dryrun,
			PSKStore:     pskStore,
			Salt:         *hostSalt,
			DeploymentId: *deploymentId,
			KeysDir:      *hostKeysPath,
			Hostnames:    hostnames,
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
		var (
			clientId     string
			clientSecret string
		)
		// Use the given value if explicitly given, otherwise take the value
		// from what should've been set at build time
		if *googleClientId != "" {
			clientId = *googleClientId
		} else {
			clientId = accord.ClientID
		}

		if *googleClientSecret != "" {
			clientSecret = *googleClientSecret
		} else {
			clientSecret = accord.ClientSecret
		}

		googleAuth := &accord.GoogleAuth{
			ClientId:      clientId,
			ClientSecret:  clientSecret,
			Domain:        *domain,
			UseWebServer:  !*nowebserver,
			WebServerPort: *webserverPort,
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
			log.Printf("remoteusername is empty, picking the same as current username: %s", usr.Username)
			remoteUsername = &usr.Username
		}
		// I don't like the pattern a lot, but I'm not sure the gain is for making a builder pattern
		// TODO: think more about this
		user := client.NewUserWithToken(c, tok)
		user.SetUsername(usr.Username)
		user.SetRemoteUsername(*remoteUsername)
		user.SetKeysDir(keysDir)
		user.SetPrincipals(principals.Value())
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
	case "trustedcerts":
		c := protocol.NewCertClient(conn)
		// Queries and prints out the trusted certs
		resp, err := c.PublicTrustedCA(context.Background(), &protocol.PublicTrustedCARequest{
			RequestTime: ptypes.TimestampNow(),
		})
		if err != nil {
			log.Fatalf("Failed to get the certs %s", err)
		}
		fmt.Println("=== Host CAs ===")
		for _, hostCA := range resp.HostCAs {
			fmt.Println(string(hostCA.GetPublicKey()))
		}

		fmt.Println("=== User CAs ===")
		for _, userCA := range resp.UserCAs {
			fmt.Println(string(userCA.GetPublicKey()))
		}
		close(done)
	case "updatehostcerts":
		c := protocol.NewCertClient(conn)
		if *knownHostsFile == "" {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			defaultPath := filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
			log.Printf("No knownhosts file given, using default: %s", defaultPath)
			knownHostsFile = &defaultPath
		}
		user := client.NewUser(c)
		err := user.UpdateHostCertAuthority(*knownHostsFile)
		if err != nil {
			log.Fatalf("Failed to update known hosts file %s. %s", *knownHostsFile, err)
		}

		close(done)
		//fmt.Printf("Resp: %#v\n", resp)
	case "updateusercerts":
		c := protocol.NewCertClient(conn)

		if *userCACertsFile == "" {
			defaultPath := "/etc/ssh/users_ca.pub"
			log.Printf("No usersCA file given, using default: %s", defaultPath)
			userCACertsFile = &defaultPath
		}
		host := client.NewHost(c)
		err := host.UpdateUserCertAuthority(*userCACertsFile)
		if err != nil {
			log.Fatalf("Failed to update trusted users ca file: %s. %s", *userCACertsFile, err)
		}
		close(done)

	case "updatesshd":
		if *sshdFile == "" {
			defaultPath := "/etc/ssh/sshd_config"
			log.Printf("No sshdConfig file given, using default: %s", defaultPath)
			sshdFile = &defaultPath
		}
		err := accord.UpdateSSHD(*sshdFile)
		if err != nil {
			log.Fatalf("Failed to write %s with %s", *sshdFile, err)
		}
		close(done)
	default:
		log.Fatalf("Don't know the task %s", *task)
	}
	<-done
	//log.Printf("Ping Response: %s. Latency: %s", r.Message, Latency(r.Metadata))
}
