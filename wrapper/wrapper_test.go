package wrapper_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/Z00mZE/pbtoken/cipher"
	"github.com/Z00mZE/pbtoken/pb"
	"github.com/Z00mZE/pbtoken/wrapper"
)

const key = "TEST_KEY_ad_mark_1234567890_____"

func TestSuccess(t *testing.T) {
	// Arrange
	remark, err := cipher.NewCipherWithNonceSize([]byte(key), 12)
	require.NoError(t, err)

	c := wrapper.NewWrapper(remark)

	ts := time.Now()
	tzName, tzOffset := ts.Zone()
	ts = ts.In(time.FixedZone(tzName, tzOffset))
	mark := &pb.Example{
		ID:          "1",
		Label:       "Label",
		Description: "description",
		Attributes: []*pb.ExampleAttribute{
			{
				ID:          "ID",
				Label:       "Label",
				Description: "Description",
				Values: []*pb.ExampleAttributeValue{
					{
						ID:          "ID",
						Label:       "Label",
						Description: "Description",
						Unit:        "mm",
						Values:      "Values",
					},
				},
			},
		},
	}

	// Action
	enc, err := c.Encode(mark)
	require.NoError(t, err)

	dec := new(pb.Example)
	err = c.Decode(enc, dec)
	require.NoError(t, err)

	// Assert
	require.Truef(t, proto.Equal(mark, dec), "not equal: \n\t\texpected:%v,\n\t\tactual:  %v", mark, dec)
}

func TestMarkPointerIsNil(t *testing.T) {
	// Arrange
	remark, err := cipher.NewCipherWithNonceSize([]byte(key), 12)
	require.NoError(t, err)
	c := wrapper.NewWrapper(remark)

	var mark *pb.Example

	// Action
	_, err = c.Encode(mark)

	// Assert
	require.Error(t, err)
}

func TestMarkIsNil(t *testing.T) {
	// Arrange
	remark, err := cipher.NewCipherWithNonceSize([]byte(key), 12)
	require.NoError(t, err)
	c := wrapper.NewWrapper(remark)

	// Action
	_, err = c.Encode(nil)

	// Assert
	require.Error(t, err)
}

func TestEncodedStringIsGarbage(t *testing.T) {
	// Arrange
	remark, err := cipher.NewCipherWithNonceSize([]byte(key), 12)
	require.NoError(t, err)
	c := wrapper.NewWrapper(remark)
	ts := time.Now()
	tzName, tzOffset := ts.Zone()
	ts = ts.In(time.FixedZone(tzName, tzOffset))
	mark := &pb.Example{
		ID:          "1",
		Label:       "Label",
		Description: "description",
		Attributes: []*pb.ExampleAttribute{
			{
				ID:          "ID",
				Label:       "Label",
				Description: "Description",
				Values: []*pb.ExampleAttributeValue{
					{
						ID:          "ID",
						Label:       "Label",
						Description: "Description",
						Unit:        "mm",
						Values:      "Values",
					},
				},
			},
		},
	}

	enc, err := c.Encode(mark)
	require.NoError(t, err)

	encNew := strings.ToUpper(enc)
	var dec pb.Example

	// Action
	err = c.Decode(encNew, &dec)

	// Assert
	require.Error(t, err)
}
