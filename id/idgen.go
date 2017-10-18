package id

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// KeyID - returns the Key ID for the server to use for connecting to the
// accprd server. Returns an error if it's not valid
// It's the top 4 bytes md5 of the salt+string written as a Big-endian uint32
// The idea is to get roughly random-looking strings that can be constructed with
// knowledge of salt and the source string, but otherwise on the wire looks random
func KeyID(s string, salt string) (uint32, error) {
	kBytes, err := KeyIDBytes(s, salt)
	if err != nil {
		return 0, fmt.Errorf("Error getting KeyID Bytes: %s", err)
	}
	return binary.BigEndian.Uint32(kBytes), nil
}

func KeyIDBytes(s string, salt string) ([]byte, error) {
	if s == "" {
		return nil, errors.New("Empty string given for ID")
	}
	h := md5.New()
	io.WriteString(h, salt)
	io.WriteString(h, s)
	sum := h.Sum(nil)
	return sum[:4], nil
}
