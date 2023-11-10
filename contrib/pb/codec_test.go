package pb

import (
	"github.com/lxzan/hasaki"
	"github.com/lxzan/hasaki/contrib/pb/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		_, err := Encoder.Encode(req)
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		_, err := Encoder.Encode(nil)
		assert.NoError(t, err)
	})

	t.Run("unexpected type", func(t *testing.T) {
		_, err := Encoder.Encode(struct{}{})
		assert.True(t, errors.Is(err, errDataType))
	})
}

func TestEncoder_ContentType(t *testing.T) {
	assert.Equal(t, Encoder.ContentType(), hasaki.MimeProtoBuf)
}

func TestDecode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		r, _ := Encoder.Encode(req)

		var res = &internal.HelloRequest{}
		var err = Decode(r, res)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, res.Name)
	})

	t.Run("unexpected type", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		r, _ := Encoder.Encode(req)

		var res = struct {
			Name string
			Age  int
		}{}
		var err = Decode(r, &res)
		assert.True(t, errors.Is(err, errDataType))
	})
}
