package accord

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"log"
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
func generateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
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

func (a *AESGCM) Encrypt(message []byte, sender uint32) ([]byte, error) {
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

	nonce, err := generateNonce(NonceSize)
	if err != nil {
		return nil, ErrEncrypt
	}

	//log.Printf("Nonce: %d %q", len(nonce), nonce)
	buf = append(buf, nonce...)
	newBuf := gcm.Seal(buf, nonce, message, buf[:4])
	return newBuf, nil
}

func (a *AESGCM) Decrypt(message []byte) ([]byte, error) {
	if len(message) <= NonceSize+4 {
		log.Println("message < noncesize + 4")
		return nil, ErrDecrypt
	}
	buf := message[:4]
	psk, err := a.store.GetPSK(message[:4])
	if err != nil {
		log.Printf("Err: failed to find the PSK for id: %q. %s", buf, err)
		return nil, ErrKeyNotFound
	}
	c, err := aes.NewCipher(psk)
	if err != nil {
		log.Printf("failed at cipher. %s", err)
		return nil, ErrDecrypt
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Printf("failed at gcm. %s", err)
		return nil, ErrDecrypt
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, message[4:])
	log.Printf("Nonce: %d %q", len(nonce), nonce)

	out, err := gcm.Open(nil, nonce, message[4+NonceSize:], message[:4])
	if err != nil {
		log.Printf("failed at open. %s", err)
		return nil, ErrDecrypt
	}

	return out, nil
}
