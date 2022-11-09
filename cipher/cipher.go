package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

// Cipher is a struct for encoding and decoding ad mark
type Cipher struct {
	aesGCM cipher.AEAD
}

// NewCipherWithNonceSize inits cipher
func NewCipherWithNonceSize(key []byte, nonceSize int) (*Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return nil, err
	}
	return &Cipher{aesGCM: aesGCM}, nil
}

// Encode encodes the mark
func (c *Cipher) Encode(data []byte) ([]byte, error) {
	nonce := make([]byte, c.aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.Wrap(err, "cannot set nonceSize")
	}
	return append(c.aesGCM.Seal(nil, nonce, data, nil), nonce...), nil
}

// Decode decodes the mark
func (c *Cipher) Decode(data []byte) ([]byte, error) {
	separator := len(data) - c.aesGCM.NonceSize()
	return c.aesGCM.Open(nil, data[separator:], data[:separator], nil)
}
