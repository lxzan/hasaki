package hasaki

import (
	"bytes"
	"encoding/xml"
	jsoniter "github.com/json-iterator/go"
	"github.com/lxzan/hasaki/internal"
	"github.com/pkg/errors"
	"github.com/valyala/bytebufferpool"
	"io"
	"net/url"
	"strings"
)

const (
	MimeJson     = "application/json;charset=utf-8"
	MimeYaml     = "application/x-yaml;charset=utf-8"
	MimeXml      = "application/xml;charset=utf-8"
	MimeProtoBuf = "application/x-protobuf"
	MimeForm     = "application/x-www-form-urlencoded"
	MimeStream   = "application/octet-stream"
	MimeJpeg     = "image/jpeg"
	MimeGif      = "image/gif"
	MimePng      = "image/png"
	MimeMp4      = "video/mpeg4"
)

type Any map[string]any

type Encoder interface {
	Encode(v any) (io.Reader, error)
	ContentType() string
}

var (
	JsonEncoder = new(jsonEncoder)
	FormEncoder = new(formEncoder)
	XmlEncoder  = new(xmlEncoder)
)

type (
	jsonEncoder struct{}
	formEncoder struct{}
	xmlEncoder  struct{}
)

func (c jsonEncoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	err := jsoniter.ConfigFastest.NewEncoder(w).Encode(v)
	r := &internal.CloserWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (c jsonEncoder) ContentType() string {
	return MimeJson
}

func JsonDecode(r io.Reader, v any) error { return jsoniter.ConfigFastest.NewDecoder(r).Decode(v) }

func (f formEncoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	switch r := v.(type) {
	case url.Values:
		return strings.NewReader(r.Encode()), nil
	case string:
		return strings.NewReader(r), nil
	default:
		return nil, errors.WithStack(errUnsupportedData)
	}
}

func (f formEncoder) ContentType() string {
	return MimeForm
}

func FormDecode(r io.Reader, v any) error {
	values, ok := v.(*url.Values)
	if !ok {
		return errors.Wrap(errUnsupportedData, "v must be *url.Values type")
	}
	var builder = &strings.Builder{}
	var buffer = internal.GetBuffer()
	var p = buffer.Bytes()[:internal.BufferSize]
	_, _ = io.CopyBuffer(builder, r, p)
	internal.PutBuffer(buffer)
	result, err := url.ParseQuery(builder.String())
	if err != nil {
		return errors.WithStack(err)
	}
	*values = result
	return nil
}

func (c xmlEncoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	err := xml.NewEncoder(w).Encode(v)
	r := &internal.CloserWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (c xmlEncoder) ContentType() string {
	return MimeXml
}

func XmlDecode(r io.Reader, v any) error {
	return errors.WithStack(xml.NewDecoder(r).Decode(v))
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
