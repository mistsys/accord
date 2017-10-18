package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/acme/autocert"

	"github.com/mistsys/accord"
	"github.com/mistsys/accord/certserver"
	"github.com/mistsys/accord/db"
	"github.com/mistsys/accord/protocol"
	"github.com/mistsys/accord/status"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	defaultPort = 50051
)

// Use actual email
const defaultContactEmail = "noreply@mistsys.com"

func GetTLS(host, cacheDir, contactEmail string) (*tls.Config, error) {
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(cacheDir),
		HostPolicy: autocert.HostWhitelist(host),
		Email:      contactEmail,
	}
	return &tls.Config{GetCertificate: manager.GetCertificate}, nil
}

func readPSKsFile(path string, psks *map[uint32][]byte) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s", path)
	}
	err = json.Unmarshal(dat, psks)
	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal json in %s", path)
	}
	return nil
}

func grpcHandlerFunc(rpcServer *grpc.Server, other http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if r.ProtoMajor == 2 && strings.Contains(ct, "application/grpc") {
			rpcServer.ServeHTTP(w, r)
		} else {
			other.ServeHTTP(w, r)
		}
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	insecure := flag.Bool("insecure", false, "Is this for development, and disable TLS?")
	// this is where reading from a HSM could be implemented
	rootCA := flag.String("rootca", "", "Path to the root CA cert. They need to be encrypted")
	userCA := flag.String("userca", "", "Path to the user CA cert. They need to be encrypted")
	port := flag.Int("port", defaultPort, "Port to use. This is overriden to 443 because letsencrypt will be used")
	psksFile := flag.String("path.psks", "", "A JSON file with all the PSKs that we're creating servers with")
	healthCheckPort := flag.Int("health.port", 9110, "Where to listen for the health checks")
	// these should only use used in insecure mode for development.. or from a config
	// read from a parameter store what the password is
	rootCAPassword := flag.String("rootcapassword", "", "The passphrase to use for the root CA")
	userCAPassword := flag.String("usercapassword", "", "The passphrase to use for the user CA")
	cacheDir := flag.String("autocert.cache", "cache", "Where to save the cached certificates, needs to be writable")
	contactEmail := flag.String("autocert.contactemail", defaultContactEmail, "Contact email to use for Lets Encrypt certificates")
	roleArn := flag.String("role-arn", "", "Role ARN to use for reading the parameter strings for root certificate")
	certsDir := flag.String("path.certs", "", "Path where certificates are -- used if role-arn is set")
	region := flag.String("aws.region", "us-east-1", "Which AWS region are we on?")
	paramsPrefix := flag.String("params-prefix", "", "Where to look for the passphrase to decrypt the HostCA and UserCA keys")
	// these should only be used for testing
	sslKey := flag.String("sslkey", "", "Path to the SSL key")
	sslCert := flag.String("sslcert", "", "Path to the SSL cert")
	hostname := flag.String("hostname", "localhost", "Hostname to use")
	// if sslcerts aren't explicity
	flag.Parse()

	var err error

	psks := make(map[uint32][]byte)
	if *psksFile == "" {
		log.Println("path.psks was empty, so initializing with default test key")
		psks[912090709] = []byte(`JpUtbRukLuIFyjeKpA4fIpjgs6MTV8eH`)
	} else { // this should be in a key in parameter store too
		err := readPSKsFile(*psksFile, &psks)
		if err != nil {
			log.Fatalf("Failed to read psk file %s. %s", *psksFile, err)
		}
	}
	// this could be loaded from another store too
	pskStore := db.NewLocalPSKStore(psks)

	var certManager *accord.CertManager

	if *roleArn != "" {
		if *certsDir == "" {
			log.Fatal("role-arn is set but certs directory isn't")
		}
		certManager, err = accord.NewCertManagerWithParameters(*certsDir, *region, *roleArn, *paramsPrefix)
		if err != nil {
			log.Fatalf("Cannot load cert manager: %s", err)
		}
	} else {
		// TODO: make it so that the certmanager scans a directory and finds IDs, then queries
		// the corresponding keys' parameters on demand. This allows us to revoke the keys as needed
		certManager, err = accord.NewCertManagerWithPasswords(*rootCA, *rootCAPassword, *userCA, *userCAPassword)
		if err != nil {
			log.Fatalf("Failed to initialize cert manager. %s", err)
		}
	}

	status.ServePort(*healthCheckPort)

	server := grpc.NewServer()
	protocol.RegisterCertServer(server, certserver.NewCertAccorder(pskStore, certManager))
	reflection.Register(server)
	addr := ":" + strconv.Itoa(*port)
	if *insecure {
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	} else {

		var tlsConfig *tls.Config
		var creds credentials.TransportCredentials
		if *sslKey != "" && *sslCert != "" {
			creds, err = credentials.NewServerTLSFromFile(*sslCert, *sslKey)
			if err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		} else {
			// Since Lets Encrypt uses SNI anyway
			// you need to run it on 443 regardless
			log.Println("Getting certificate from Lets Encrypt")
			addr = ":https"
			tlsConfig, err = GetTLS(*hostname, *cacheDir, *contactEmail)
			if err != nil {
				log.Fatalf("Failed to get certificate for %s", *hostname)
			}
			creds = credentials.NewTLS(tlsConfig)
		}

		// TODO: refactor this, check for certs first or use config
		server := grpc.NewServer(grpc.Creds(creds))

		protocol.RegisterCertServer(server, certserver.NewCertAccorder(pskStore, certManager))
		reflection.Register(server)
		mux := http.DefaultServeMux
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
		})
		srv := &http.Server{
			Addr:      addr,
			Handler:   grpcHandlerFunc(server, mux),
			TLSConfig: tlsConfig,
		}

		if err := srv.ListenAndServeTLS(*sslCert, *sslKey); err != nil {
			log.Fatalf("Error running server %s", err)
		}

		//lis, err := net.Listen("tcp", addr)
		//if err != nil {
		//	log.Fatalf("Failed to listen: %v", err)
		//}
		//if err := server.Serve(lis); err != nil {
		//	log.Fatalf("Failed to serve: %v", err)
		//}
	}
}
