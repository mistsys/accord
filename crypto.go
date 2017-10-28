package accord

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"
	"log"

	"github.com/pkg/errors"
)

const (
	KeySize   = 32
	NonceSize = 12
)

var (
	ErrEncrypt     = errors.New("secret: encryption failed")
	ErrDecrypt     = errors.New("secret: decryption failed")
	ErrKeyNotFound = errors.New("keylookup: key not found")
)

// GenerateNonce creates a new random nonce.
func GenerateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// len(encodeURL) == 64. This allows (x <= 265) x % 64 to have an even
// distribution.
const urlEncodeable = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// A helper function create and fill a slice of length n with characters from
// a-zA-Z0-9_-. It panics if there are any problems getting random bytes.
func RandAsciiBytes(n int) []byte {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(len(urlEncodeable))

		// put into output
		output[pos] = urlEncodeable[randomPos]
	}

	return output
}

// GenerateKey generates a new AES-256 key.
func GenerateKey() []byte {
	return RandAsciiBytes(KeySize)
}

// Initialize the AESGCM with a PSK store, this can be anything from a local instance
// or something that reads from a HSM or memory, the logic for getting the key securely
// will be in the PSKStore implmentation
type AESGCM struct {
	store PSKStore
}

func InitAESGCM(store PSKStore) *AESGCM {
	return &AESGCM{
		store: store,
	}
}

func (a *AESGCM) EncryptWithNonce(message []byte, nonce []byte, sender uint32) ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, sender)

	psk, err := a.store.GetPSK(buf)
	if err != nil {
		log.Printf("Err: failed to find the PSK for id: %d. %s", sender, err)
		return nil, ErrKeyNotFound
	}
	c, err := aes.NewCipher(psk)
	if err != nil {
		log.Printf("Err: %s", err)
		return nil, ErrEncrypt
	}

	//log.Printf("Cipher Block Size: %d", c.BlockSize())

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, ErrEncrypt
	}

	//log.Printf("Nonce: %d %q", len(nonce), nonce)
	buf = append(buf, nonce...)
	newBuf := gcm.Seal(buf, nonce, message, buf[:4])
	return newBuf, nil
}

func (a *AESGCM) Encrypt(message []byte, sender uint32) ([]byte, error) {
	nonce, err := GenerateNonce(NonceSize)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to generate Nonce")
	}
	return a.EncryptWithNonce(message, nonce, sender)
}

// Decrypt returns plaintext, nonce, senderid, error
func (a *AESGCM) Decrypt(message []byte) ([]byte, []byte, uint32, error) {
	if len(message) <= NonceSize+4 {
		log.Println("message < noncesize + 4")
		return nil, nil, 0, ErrDecrypt
	}
	buf := message[:4]
	sender := binary.BigEndian.Uint32(buf)
	psk, err := a.store.GetPSK(message[:4])
	if err != nil {
		log.Printf("Err: failed to find the PSK for id: %q. %s", buf, err)
		return nil, nil, sender, ErrKeyNotFound
	}
	c, err := aes.NewCipher(psk)
	if err != nil {
		log.Printf("failed at cipher. %s", err)
		return nil, nil, sender, ErrDecrypt
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Printf("failed at gcm. %s", err)
		return nil, nil, sender, ErrDecrypt
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, message[4:])
	log.Printf("Nonce: %d %q", len(nonce), nonce)

	out, err := gcm.Open(nil, nonce, message[4+NonceSize:], message[:4])
	if err != nil {
		log.Printf("failed at open. %s", err)
		return nil, nil, sender, ErrDecrypt
	}

	return out, nonce, sender, nil
}
