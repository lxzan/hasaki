//go:generate go mod tidy

package pb

import (
	"bytes"
	"github.com/lxzan/hasaki"
	"github.com/lxzan/hasaki/internal"
	"github.com/pkg/errors"
	"github.com/valyala/bytebufferpool"
	"google.golang.org/protobuf/proto"
	"io"
)

var (
	errDataType = errors.New("v must be proto.Message type")
	Encoder     = new(encoder)
)

type encoder struct{}

func (c encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	message, ok := v.(proto.Message)
	if !ok {
		return nil, errors.WithStack(errDataType)
	}
	p, err := proto.Marshal(message)
	return bytes.NewReader(p), errors.WithStack(err)
}

func (c encoder) ContentType() string {
	return hasaki.MimeProtoBuf
}

func Decode(r io.Reader, v any) error {
	message, ok := v.(proto.Message)
	if !ok {
		return errors.WithStack(errDataType)
	}
	var w = bytebufferpool.Get()
	var buffer = internal.GetBuffer()
	var p = buffer.Bytes()[:internal.BufferSize]
	_, _ = io.CopyBuffer(w, r, p)
	internal.PutBuffer(buffer)
	bytebufferpool.Put(w)
	return errors.WithStack(proto.Unmarshal(w.B, message))
}
