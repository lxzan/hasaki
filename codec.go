package hasaki

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/url"
	"strings"

	"github.com/valyala/bytebufferpool"

	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	MimeJson   = "application/json;charset=utf-8"
	MimeForm   = "application/x-www-form-urlencoded"
	MimeStream = "application/octet-stream"
	MimeJpeg   = "image/jpeg"
	MimeGif    = "image/gif"
	MimePng    = "image/png"
	MimeMp4    = "video/mpeg4"
)

type Any map[string]any

type Encoder interface {
	Encode(v any) (io.Reader, error)
	ContentType() string
}

var (
	JsonEncoder = new(json_encoder)
	FormEncoder = new(form_encoder)
)

type json_encoder struct{}

func (j json_encoder) Encode(v any) (io.Reader, error) {
	w := bytebufferpool.Get()
	err := jsoniter.ConfigFastest.NewEncoder(w).Encode(v)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (j json_encoder) ContentType() string {
	return MimeJson
}

type form_encoder struct{}

// Encode do not support float number
func (f form_encoder) Encode(v any) (io.Reader, error) {
	if values, ok := v.(url.Values); ok {
		return strings.NewReader(values.Encode()), nil
	}
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return strings.NewReader(values.Encode()), nil
}

func (f form_encoder) ContentType() string {
	return MimeForm
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

func JsonDecode(r io.Reader, v any) error { return jsoniter.ConfigFastest.NewDecoder(r).Decode(v) }

func XmlDecode(r io.Reader, v any) error { return xml.NewDecoder(r).Decode(v) }
