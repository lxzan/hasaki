package pb

import (
	"bytes"
	"github.com/lxzan/hasaki"
	"github.com/lxzan/hasaki/contrib/pb/internal"
	internal2 "github.com/lxzan/hasaki/internal"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/bytebufferpool"
	"io"
	"testing"
)

func TestEncoder_Encode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		_, err := Codec.Encode(req)
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		_, err := Codec.Encode(nil)
		assert.NoError(t, err)
	})

	t.Run("unexpected type", func(t *testing.T) {
		_, err := Codec.Encode(struct{}{})
		assert.True(t, errors.Is(err, errDataType))
	})
}

func TestEncoder_ContentType(t *testing.T) {
	assert.Equal(t, Codec.ContentType(), hasaki.MimeProtoBuf)
}

func TestDecode(t *testing.T) {
	t.Run("ok 1", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		r, _ := Codec.Encode(req)

		var res = &internal.HelloRequest{}
		var err = Codec.Decode(r, res)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, res.Name)
	})

	t.Run("ok 2", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		r, _ := Codec.Encode(req)
		p, _ := io.ReadAll(r)
		br := &internal2.CloserWrapper{B: &bytebufferpool.ByteBuffer{B: p}, R: bytes.NewReader(p)}

		var res = &internal.HelloRequest{}
		var err = Codec.Decode(br, res)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, res.Name)
	})

	t.Run("unexpected type", func(t *testing.T) {
		var req = &internal.HelloRequest{
			Name: "caster",
			Age:  1,
		}
		r, _ := Codec.Encode(req)

		var res = struct {
			Name string
			Age  int
		}{}
		var err = Codec.Decode(r, &res)
		assert.True(t, errors.Is(err, errDataType))
	})
}
