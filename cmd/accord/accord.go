package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

/*
I wrote this to understand what was going on when a SSH certificate is generated, managed
This is left here as an example and possibly as a future user of the underlying API
This tool is quite sufficient if you like the idea of running your own pure Go SSH CA service
or need to wrap it into another tool or language that does the same thing
*/

func parsePubKey(pubKeyPath string) {
	contents, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Failed to read file %s. %s", pubKeyPath, err)
	}

	_, comment, options, rest, err := ssh.ParseAuthorizedKey(contents)
	if err != nil {
		log.Fatalf("Failed to parse cert key: %s. %s", string(contents), err)
	}

	fmt.Printf("Parsed: %s\n", pubKeyPath)
	fmt.Println(comment, options, string(rest))
}

func MarshalCert(cert *ssh.Certificate, comment string) []byte {
	b := &bytes.Buffer{}
	b.WriteString(cert.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(cert.Marshal())
	e.Close()
	b.WriteByte(' ')
	b.WriteString(comment)
	b.WriteByte('\n')
	return b.Bytes()
}

// ssh one-line format (for lack of a better term) consists of three text fields: { key_type, data, comment }
// data is base64 encoded binary which consists of tuples of length (4 bytes) and data of the length described previously.
// For RSA keys, there should be three tuples which should be:  { key_type, public_exponent, modulus }
func EncodeRSAPublicKey(key interface{}, comment string) (string, error) {
	if rsaKey, ok := key.(rsa.PublicKey); ok {
		key_type := "ssh-rsa"

		modulus_bytes := rsaKey.N.Bytes()

		buf := new(bytes.Buffer)

		var data = []interface{}{
			uint32(len(key_type)),
			[]byte(key_type),
			uint32(binary.Size(uint32(rsaKey.E))),
			uint32(rsaKey.E),
			uint32(binary.Size(modulus_bytes)),
			modulus_bytes,
		}

		for _, v := range data {
			err := binary.Write(buf, binary.BigEndian, v)
			if err != nil {
				return "", err
			}
		}

		return fmt.Sprintf("%s %s %s", key_type, base64.StdEncoding.EncodeToString(buf.Bytes()), comment), nil
	}

	return "", fmt.Errorf("Unknown key type: %T\n", key)
}

func main() {
	certKeyPath := flag.String("certkey", "", "Path for the certificate to use for signing")
	pubKeyPath := flag.String("pubkey", "", "SSH Public Key to sign with the cert")
	pubCertPath := flag.String("pubcert", "", "Generated Cert Path")
	hostname := flag.String("hostname", "", "Hostname to sign the cert for")
	password := flag.String("password", "", "Password to encrypt the root key with")
	task := flag.String("task", "genusercert", "Task to do")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch *task {
	case "genusercert":
		var (
			err    error
			signer ssh.Signer
		)

		contents, err := ioutil.ReadFile(*certKeyPath)
		if err != nil {
			log.Fatalf("Failed to read file %s. %s", *certKeyPath, err)
		}

		if *password != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(contents, []byte(*password))
			if err != nil {
				log.Fatalf("Failed to parse private key: %s", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(contents)
			if err != nil {
				log.Fatalf("Failed to parse private key: %s", err)
			}
		}

		contents, err = ioutil.ReadFile(*pubKeyPath)
		if err != nil {
			log.Fatalf("Failed to read file %s. %s", *pubKeyPath, err)
		}

		pubkey, comment, _, _, err := ssh.ParseAuthorizedKey(contents)
		if err != nil {
			log.Fatalf("Failed to parse pub key: %s. %s", *pubKeyPath, err)
		}

		curTime := time.Now()
		oneDayLater := curTime.Add(24 * time.Hour)

		cert := &ssh.Certificate{
			CertType:        ssh.UserCert,
			Key:             pubkey,
			KeyId:           os.Getenv("USER"),
			Serial:          1,
			ValidBefore:     uint64(oneDayLater.Unix()),
			ValidAfter:      uint64(curTime.Unix()),
			ValidPrincipals: []string{"pgautam", "admin"},
		}
		err = cert.SignCert(rand.Reader, signer)
		if err != nil {
			log.Fatalf("failed to sign the cert. %s", err)
		}
		fmt.Printf(string(MarshalCert(cert, comment)))

	// this is test code to see if go can generate valid rsa cert key
	// it should be equivalent to
	// ssh-keygen -t rsa -b 4096 -C "Cert valid from <YYYY-MM-DD> to <YYYY-MM-DD+90>"
	case "genrootcert":
		argv := flag.Args()
		wPrivKey := os.Stdout
		wPubKey := os.Stdout
		log.Printf("len(argv) = %d", len(argv))
		var err error
		if len(argv) == 1 {
			filename := argv[0]
			wPrivKey, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)

			if err != nil {
				log.Fatalf("Failed to create %s for private key", filename)
			}
			defer wPrivKey.Close()

			pubKeyFilename := filename + ".pub"
			wPubKey, err = os.Create(pubKeyFilename)

			if err != nil {
				log.Fatalf("Failed to create %s for public key", pubKeyFilename)
			}
			defer wPubKey.Close()

		}
		curTime := time.Now()
		ninetyDaysLater := curTime.Add(90 * 24 * time.Hour)
		privkey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatalf("Failed to generate RSA key. %s", err)
		}

		privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
		block := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		}

		if *password != "" {
			block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, []byte(*password), x509.PEMCipherAES256)
			if err != nil {
				log.Fatalf("Failed to encrypt PEM block with the provided password %s", err)
			}
		}

		fmt.Println("==== PRIVATE KEY ====")
		err = pem.Encode(wPrivKey, block)
		if err != nil {
			log.Fatalf("encode private key failed")
		}
		// parse DER format to a native type
		key, err := x509.ParsePKCS1PrivateKey(privkey_bytes)
		if err != nil {
			log.Fatal("Failed to parse PKCS1 Private Key")
		}

		// encode the public key portion of the native key into ssh-rsa format
		// second parameter is the optional "comment" at the end of the string (usually 'user@host')
		ssh_rsa, err := EncodeRSAPublicKey(key.PublicKey, fmt.Sprintf("Cert valid from %s to %s", curTime.Format("2006-01-02"), ninetyDaysLater.Format("2006-01-02")))
		if err != nil {
			log.Fatal("Failed to encode RSA Public Key")
		}

		fmt.Println("==== PUBLIC KEY ====")
		fmt.Fprintf(wPubKey, ssh_rsa)

		//fmt.Printf(string(MarshalCert(cert, comment)))
		//fmt.Println("Public Key")

		//fmt.Printf(stri)

	case "genhostcert":
		if *hostname == "" {
			log.Fatal("-hostname is required")
		}

		var (
			signer ssh.Signer
			err    error
		)

		contents, err := ioutil.ReadFile(*certKeyPath)
		if err != nil {
			log.Fatalf("Failed to read file %s. %s", *certKeyPath, err)
		}

		if *password != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(contents, []byte(*password))
			if err != nil {
				log.Fatalf("Failed to parse private key: %s", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKey(contents)
			if err != nil {
				log.Fatalf("Failed to parse private key: %s", err)
			}
		}

		contents, err = ioutil.ReadFile(*pubKeyPath)
		if err != nil {
			log.Fatalf("Failed to read file %s. %s", *pubKeyPath, err)
		}

		pubkey, comment, _, _, err := ssh.ParseAuthorizedKey(contents)
		if err != nil {
			log.Fatalf("Failed to parse pub key: %s. %s", *pubKeyPath, err)
		}

		curTime := time.Now()
		oneDayLater := curTime.Add(24 * time.Hour)

		cert := &ssh.Certificate{
			CertType:        ssh.HostCert,
			Key:             pubkey,
			Serial:          1,
			ValidBefore:     uint64(oneDayLater.Unix()),
			ValidAfter:      uint64(curTime.Unix()),
			ValidPrincipals: []string{*hostname},
		}
		err = cert.SignCert(rand.Reader, signer)
		if err != nil {
			log.Fatalf("failed to sign the cert. %s", err)
		}
		fmt.Printf(string(MarshalCert(cert, comment)))
	case "printcert":
		contents, err := ioutil.ReadFile(*pubCertPath)
		if err != nil {
			log.Fatalf("Failed to read file %s. %s", *pubCertPath, err)
		}
		key, _, _, _, err := ssh.ParseAuthorizedKey(contents)
		cert, ok := key.(*ssh.Certificate)
		if !ok {
			log.Fatalf("got %v (%T), wanted a certificate")
		}
		fmt.Printf("%#v\n", cert)
		//fmt.Printf("Signature: %#v\n", cert.Signature)
		//fmt.Printf(string(MarshalCert(cert)))
		//fmt.Println(comment)
	}

}
