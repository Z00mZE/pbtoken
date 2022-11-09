package cipher

import (
	"crypto/rand"
	"io"
	"reflect"
	"testing"
)

func TestNewCipherWithNonceSize(t *testing.T) {
	type args struct {
		key       []byte
		nonceSize int
	}
	type testPlan struct {
		name    string
		args    args
		wantErr bool
	}
	short := testPlan{
		name: "Short",
		args: args{
			key:       []byte(`TEST_KEY_ad_mark_1234567890____`),
			nonceSize: 0,
		},
		wantErr: true,
	}
	normal := testPlan{
		name: "Normal",
		args: args{
			key:       []byte(`TEST_KEY_ad_mark_1234567890_____`),
			nonceSize: 1,
		},
		wantErr: false,
	}
	long := testPlan{
		name: "long",
		args: args{
			key:       []byte(`TEST_KEY_ad_mark_1234567890______`),
			nonceSize: 1,
		},
		wantErr: true,
	}
	tests := []testPlan{short, normal, long}

	for i := 0; i < 32; i++ {
		short.args.nonceSize *= 2
		normal.args.nonceSize *= 2
		long.args.nonceSize *= 2
		tests = append(
			tests,
			short,
			normal,
			long,
		)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCipherWithNonceSize(tt.args.key, tt.args.nonceSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCipherWithNonceSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCipher_EncodeDecode(t *testing.T) {
	type testCause struct {
		name     string
		key      []byte
		nounSize int
	}

	keygen := func(length int) []byte {
		out := make([]byte, length)
		_, _ = io.ReadFull(rand.Reader, out)
		return out
	}
	c, cError := NewCipherWithNonceSize(keygen(32), 12)
	if cError != nil {
		t.Errorf("NewCipherWithNonceSize() error = %v;", cError)
		return
	}
	for i := 0; i < 1_000; i++ {
		msg := keygen(i)
		enc, encError := c.Encode(msg)
		if encError != nil {
			t.Errorf("Encode() error = %v", encError)
			return
		}
		dec, msgError := c.Decode(enc)
		if msgError != nil {
			t.Errorf("Encode() error = %v", encError)
			return
		}
		if !reflect.DeepEqual(string(msg), string(dec)) {
			t.Errorf("Encode->Decode() got = [%v], want [%v]", string(dec), string(msg))
		}
	}
}
