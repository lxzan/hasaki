package hasaki

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"github.com/valyala/bytebufferpool"

	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	MimeJSON   = "application/json;charset=utf-8"
	MimeFORM   = "application/x-www-form-urlencoded"
	MimeSTREAM = "application/octet-stream"
	MimeJPEG   = "image/jpeg"
	MimeGIF    = "image/gif"
	MimePNG    = "image/png"
	MimeMP4    = "video/mpeg4"
)

type Any map[string]any

type Encoder interface {
	Encode(v any) (io.Reader, error)
	ContentType() string
}

var (
	JSONEncoder = new(jsonEncoder)
	FORMEncoder = new(formEncoder)
)

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v any) (io.Reader, error) {
	w := bytebufferpool.Get()
	err := jsoniter.ConfigFastest.NewEncoder(w).Encode(v)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (j jsonEncoder) ContentType() string {
	return MimeJSON
}

type formEncoder struct{}

// Encode do not support float number
func (f formEncoder) Encode(v any) (io.Reader, error) {
	if values, ok := v.(url.Values); ok {
		return strings.NewReader(values.Encode()), nil
	}
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return strings.NewReader(values.Encode()), nil
}

func (f formEncoder) ContentType() string {
	return MimeFORM
}

type closerWrapper struct {
	B *bytebufferpool.ByteBuffer
	R io.Reader
}

func (c *closerWrapper) Read(p []byte) (n int, err error) {
	return c.R.Read(p)
}

func (c *closerWrapper) Close() error {
	bytebufferpool.Put(c.B)
	return nil
}

type streamEncoder struct {
	contentType string
}

func NewStreamEncoder(contentType string) Encoder {
	return &streamEncoder{contentType: contentType}
}

func (c *streamEncoder) Encode(v any) (io.Reader, error) {
	switch r := v.(type) {
	case io.Reader:
		return r, nil
	case []byte:
		return bytes.NewReader(r), nil
	case string:
		return strings.NewReader(r), nil
	default:
		return nil, errors.WithStack(errUnsupportedData)
	}
}

func (c *streamEncoder) ContentType() string {
	return c.contentType
}
