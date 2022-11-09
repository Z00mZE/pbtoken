package wrapper

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const nonceSize = 12

type Cipher interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}

// Wrapper is a wrapper over cipher
type Wrapper struct {
	cipher     Cipher
	serializer *base64.Encoding
}

// NewWrapper inits new Cipher with package nonce
func NewWrapper(cipher Cipher) *Wrapper {
	return &Wrapper{cipher: cipher, serializer: base64.RawURLEncoding.Strict()}
}

// Encode implements encoding
func (w *Wrapper) Encode(msg proto.Message) (string, error) {
	if msg == nil || !msg.ProtoReflect().IsValid() {
		return "", errors.New("invalid data")
	}
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return "", errors.Wrap(err, "marshaling proto.Message failed: error")
	}
	encrypted, outError := w.cipher.Encode(bytes)
	if outError != nil {
		return "", outError
	}
	return w.serializer.EncodeToString(encrypted), nil
}

// Decode implements decoding
func (w *Wrapper) Decode(data string, msg proto.Message) error {
	if msg == nil || !msg.ProtoReflect().IsValid() {
		return errors.New("invalid receiver")
	}

	decoded, decodedError := w.serializer.DecodeString(data)
	if decodedError != nil {
		return errors.New("Unmarshalling base64 failed")
	}
	decoded, deccodedError := w.cipher.Decode(decoded)
	if deccodedError != nil {
		return errors.Wrap(deccodedError, "wrapper.cipher.Decode(): error")
	}

	return proto.Unmarshal(decoded, msg)
}
