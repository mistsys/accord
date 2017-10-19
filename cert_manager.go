package accord

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/mistsys/accord/aws_params"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	ErrInvalidSerial      = errors.New("Serial should be something meaningful, not 0")
	ErrInvalidStartTime   = errors.New("Cannot sign for certs with time in the past")
	ErrEndBeforeStartTime = errors.New("End Time cannot be before start time")
	ErrEmptyID            = errors.New("Empty ID supplied")
	ErrValidityTooLong    = errors.New("The Validity for certs is too long")
	keyPairRegex          = regexp.MustCompile(`ca_(?P<type>user|host)_(?P<id>\d+).?(?P<key_type>pub)?`)
)

// TODO: this needs a real refactoring
// I designed it to be very simple and it grew relatively organic
// Once the interfaces are little cleaner, update the CertManager to not store everything

// I have doubts about this approach
// the passwords are stored in plaintext in memory
// but the actual certs are read on demand, and not kept in memory
// so any kind of overflow attack needs a filesystem access too
// while this may leak the CA Passwords, as long as the certs are
// read on demand, and forgotten immediately, I think this should be
// relatively safe. This can be an interface to read from different
// stores for secrets, but that needs a lot more testing right now
// TODO: refactor to use CACertPair structure
type CertManager struct {
	rootCAPath       string
	rootCAId         int
	rootCAPassword   string
	rootCAPubKey     ssh.PublicKey
	rootCAValidFrom  time.Time
	rootCAValidUntil time.Time
	userCAId         int
	userCAPath       string
	userCAPassword   string
	userCAPubKey     ssh.PublicKey
	userCAValidFrom  time.Time
	userCAValidUntil time.Time
}

// CertMetadata is included in the comment for public key
type CertMetadata struct {
	Id         int       `json:"id"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
}

type CAType string

const (
	User CAType = "user"
	Host CAType = "host"
)

type CACertPair struct {
	Type           CAType
	Metadata       CertMetadata
	PrivateKeyPath string
	PublicKey      ssh.PublicKey
	PublicKeyPath  string
}

func (c *CACertPair) updateMetadata() error {
	var metadata = CertMetadata{}
	contents, err := ioutil.ReadFile(c.PublicKeyPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read %s", c.PublicKeyPath)
	}
	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(contents)
	if err != nil {
		return errors.Wrapf(err, "%s doesn't look like a public key file", c.PublicKeyPath)
	}

	err = json.Unmarshal([]byte(comment), &metadata)
	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal comment")
	}
	c.PublicKey = pubKey
	c.Metadata = metadata
	return nil
}

func matchFileNames(name string) map[string]string {
	matches := keyPairRegex.FindAllStringSubmatch(name, -1)
	if len(matches) == 0 {
		return nil
	}
	subExps := keyPairRegex.SubexpNames()
	elements := matches[0]
	md := map[string]string{}
	for i := 1; i < len(elements); i++ {
		md[subExps[i]] = elements[i]
	}
	return md
}

// Returns full path for all the keys in the directory that are in the following format
// ca_(user|server)_<id>
// both private key and public keys are scanned for
func certPairsInDir(dir string) (map[int]*CACertPair, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to enumerate files from %s", dir)
	}
	pairs := make(map[int]*CACertPair)
	var pair *CACertPair
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		matches := matchFileNames(name)
		if matches == nil {
			continue
		}
		idStr := matches["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to read %s, invalid id in filename", name)
		}
		if _, ok := pairs[id]; !ok {
			pairs[id] = &CACertPair{}
		}

		pair = pairs[id]
		if matches["key_type"] != "pub" {
			pair.PrivateKeyPath = filepath.Join(dir, name)
		} else {
			pair.PublicKeyPath = filepath.Join(dir, name)
			err := pair.updateMetadata()
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse public key %s for metadata", name)
			}
		}

		if matches["type"] == "host" {
			pair.Type = Host
		} else {
			pair.Type = User
		}
	}
	return pairs, nil
}

// getPassphrase queries the parameter store key at `{paramsPrefix}/id` and gets the passphrase`
func getPassphrase(client aws_params.Client, paramsPrefix string, id int) (string, error) {
	path := paramsPrefix + "/" + strconv.Itoa(id)
	return client.GetSecureString(path)
}

// NewCertManagerwithParameters looks for files ending ca_(user|host)_(identifier) and ca_(user|host)_(identifier).pub
// in the certsDirectory, the comment field is read from the corresponding public key file
// the comment in the public key file contains how long the certificate is valid for
// the identifier is used to lookup the passphrase in the parameter store with _(identifier) appended to paramsPrefix
// TODO: this should be implementing an interface so that this file doesn't remain dirty with aws deps
func NewCertManagerWithParameters(certsDir string, region string, roleArn string, paramsPrefix string) (*CertManager, error) {
	config := aws_params.NewConfig(region)
	config.RoleArn = roleArn
	client, err := aws_params.NewClient(config)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot create new client")
	}
	certPairs, err := certPairsInDir(certsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find cert pairs in %s", certsDir)
	}
	certManager := &CertManager{}
	for id, certPair := range certPairs {
		switch certPair.Type {
		case User:
			certManager.userCAPath = certPair.PrivateKeyPath
			passphrase, err := getPassphrase(client, paramsPrefix, id)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to read the user key for id %d", id)
			}
			certManager.userCAId = certPair.Metadata.Id
			certManager.userCAPassword = passphrase
			certManager.userCAPubKey = certPair.PublicKey
			certManager.userCAValidFrom = certPair.Metadata.ValidFrom
			certManager.userCAValidUntil = certPair.Metadata.ValidUntil
		case Host:
			certManager.rootCAPath = certPair.PrivateKeyPath
			passphrase, err := getPassphrase(client, paramsPrefix, id)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to read the host key for id %d", id)
			}
			certManager.rootCAId = certPair.Metadata.Id
			certManager.rootCAPassword = passphrase
			certManager.rootCAPubKey = certPair.PublicKey
			certManager.rootCAValidFrom = certPair.Metadata.ValidFrom
			certManager.rootCAValidUntil = certPair.Metadata.ValidUntil
		}
	}

	return certManager, nil
}

// NewCertmanagerwithPasswords just reads the files and decrypts them with corresponding passwords
func NewCertManagerWithPasswords(rootCAPath string, rootCAPassword string,
	userCAPath string, userCAPassword string) (*CertManager, error) {
	// Read the keys at initialization time and keep them
	// so that every time a public key is requested, we're not reading from the
	// file and decrypting the passphrases
	rootCASigner, err := getSigner(rootCAPath, rootCAPassword)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot read rootCA from %s", rootCAPath)
	}
	userCASigner, err := getSigner(userCAPath, userCAPassword)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot read userCA from %s", userCAPath)
	}

	now := time.Now()
	ninetyDaysLater := now.Add(24 * 90 * time.Hour)

	return &CertManager{
		rootCAPath:       rootCAPath,
		rootCAPassword:   rootCAPassword,
		rootCAPubKey:     rootCASigner.PublicKey(),
		rootCAValidFrom:  now,
		rootCAValidUntil: ninetyDaysLater,
		userCAPath:       userCAPath,
		userCAPassword:   userCAPassword,
		userCAPubKey:     userCASigner.PublicKey(),
		userCAValidFrom:  now,
		userCAValidUntil: ninetyDaysLater,
	}, nil
}

// ssh one-line format (for lack of a better term) consists of three text fields: { key_type, data, comment }
// data is base64 encoded binary which consists of tuples of length (4 bytes) and data of the length described previously.
// For RSA keys, there should be three tuples which should be:  { key_type, public_exponent, modulus }
func EncodePublicKey(key interface{}, comment string) (string, error) {
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

// CertSignRequest is intended to be usable from other go libraries
// easily with the start and end times to be in time.Time. Since the upstream
// service is likely only seeing the bytes for a public key, the parsing is
// done here too
type CertSignRequest struct {
	PubKey     []byte
	ValidFrom  time.Time
	ValidUntil time.Time
	Id         string
	Serial     uint64
	Principals []string
	// include the criticalOptions and Extensions too
	ssh.Permissions
}

func (r *CertSignRequest) valid() (bool, error) {
	if r.Serial == 0 {
		return false, ErrInvalidSerial
	}
	if r.ValidFrom.Before(time.Now()) {
		return false, ErrInvalidStartTime
	}
	if r.ValidUntil.Before(r.ValidFrom) {
		return false, ErrEndBeforeStartTime
	}

	if r.Id == "" {
		return false, ErrEmptyID
	}

	// keep some sane bounds on how long a cert can be for
	// 90 days is probably too long
	if r.ValidUntil.Sub(r.ValidFrom) > 90*24*time.Hour {
		return false, ErrValidityTooLong
	}

	return true, nil
}

// only allowing path so that it can be handled accordingly
func getSigner(certKeyPath, passphrase string) (ssh.Signer, error) {
	var (
		err    error
		signer ssh.Signer
	)
	contents, err := ioutil.ReadFile(certKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read the file: %s", certKeyPath)
	}
	// this should be the happy path
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(contents, []byte(passphrase))
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to parse key file: %s", certKeyPath)
		}

	} else {
		signer, err = ssh.ParsePrivateKey(contents)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to parse key file: %s", certKeyPath)
		}
	}
	return signer, err
}

func (m *CertManager) rootCASigner() (ssh.Signer, error) {
	return getSigner(m.rootCAPath, m.rootCAPassword)
}

func (m *CertManager) userCASigner() (ssh.Signer, error) {
	return getSigner(m.userCAPath, m.userCAPassword)
}

func (m *CertManager) marshalCert(cert *ssh.Certificate, comment string) []byte {
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

// these return array because we have to overlap multiple root
// and user keys when we need to rotate the keys in future
// it makes sense to make the API exposed to users be a little more flexible
func (m *CertManager) RootCAPublicKeys() []ssh.PublicKey {
	return []ssh.PublicKey{m.rootCAPubKey}
}

// CAPublic is the public data that we want to return the user
type CAPublic struct {
	Id         int       `json:"id"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
	PublicKey  []byte
}

func (m *CertManager) HostCAs() []CAPublic {
	return []CAPublic{{
		Id:         m.rootCAId,
		PublicKey:  ssh.MarshalAuthorizedKey(m.rootCAPubKey),
		ValidFrom:  m.rootCAValidFrom,
		ValidUntil: m.rootCAValidUntil,
	}}
}

func (m *CertManager) UserCAs() []CAPublic {
	return []CAPublic{{
		Id:         m.userCAId,
		PublicKey:  ssh.MarshalAuthorizedKey(m.userCAPubKey),
		ValidFrom:  m.userCAValidFrom,
		ValidUntil: m.userCAValidUntil,
	}}
}

func (m *CertManager) UserCAPublicKeys() []ssh.PublicKey {
	return []ssh.PublicKey{m.userCAPubKey}
}

// JoinPublickeys encodes the public key and joins them in a single bytearray
func JoinPublicKeys(keys []ssh.PublicKey) []byte {
	b := &bytes.Buffer{}
	for _, k := range keys {
		encoded := k.Marshal()
		b.Write(encoded)
		b.WriteByte('\n')
	}
	return b.Bytes()
}

func (m *CertManager) SignUserCert(request *CertSignRequest) ([]byte, error) {
	// Don't bother if the request is invalid
	_, err := request.valid()
	if err != nil {
		return nil, errors.Wrapf(err, "CertSignRequest isn't valid")
	}

	pubkey, comment, _, _, err := ssh.ParseAuthorizedKey(request.PubKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse public key")
	}

	cert := &ssh.Certificate{
		CertType:        ssh.UserCert,
		Key:             pubkey,
		KeyId:           request.Id,
		Serial:          request.Serial,
		ValidAfter:      uint64(request.ValidFrom.Unix()),
		ValidBefore:     uint64(request.ValidUntil.Unix()),
		ValidPrincipals: request.Principals,
	}

	signer, err := m.userCASigner()
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot get the User CA signer")
	}
	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sign the cert")
	}
	b := m.marshalCert(cert, comment)
	return b, nil
}

func (m *CertManager) SignHostCert(request *CertSignRequest) ([]byte, error) {
	// Don't bother if the request is invalid
	_, err := request.valid()
	if err != nil {
		return nil, errors.Wrapf(err, "CertSignRequest isn't valid")
	}
	pubkey, comment, _, _, err := ssh.ParseAuthorizedKey(request.PubKey)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse public key")
	}

	//log.Println("Principals ", request.)
	cert := &ssh.Certificate{
		CertType:        ssh.HostCert,
		Key:             pubkey,
		KeyId:           request.Id,
		Serial:          request.Serial,
		ValidAfter:      uint64(request.ValidFrom.Unix()),
		ValidBefore:     uint64(request.ValidUntil.Unix()),
		ValidPrincipals: request.Principals,
	}

	signer, err := m.rootCASigner()
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot get the Host CA signer")
	}
	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to sign the cert")
	}
	b := m.marshalCert(cert, comment)
	return b, nil
}
