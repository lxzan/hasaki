package hasaki

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/url"
	"strings"

	"github.com/valyala/bytebufferpool"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"

	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	MimeJson   = "application/json;charset=utf-8"
	MimeYaml   = "application/x-yaml;charset=utf-8"
	MimeXml    = "application/xml;charset=utf-8"
	MimeProto  = "application/x-protobuf"
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

type BytesCodec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var (
	JsonEncoder  = new(json_encoder)
	XmlEncoder   = new(xml_encoder)
	YamlEncoder  = new(yaml_encoder)
	ProtoEncoder = new(proto_encoder)
	FormEncoder  = new(form_encoder)
)

var (
	JSONCodec  BytesCodec = jsoniter.ConfigCompatibleWithStandardLibrary
	YAMLCodec  BytesCodec = new(yaml_codec)
	XMLCodec   BytesCodec = new(xml_codec)
	PROTOCodec BytesCodec = new(proto_codec)
)

type json_encoder struct{}

func (j json_encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	b, err := JSONCodec.Marshal(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	w.Set(b)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (j json_encoder) ContentType() string {
	return MimeJson
}

type yaml_encoder struct{}

func (y yaml_encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	b, err := YAMLCodec.Marshal(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	w.Set(b)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (y yaml_encoder) ContentType() string {
	return MimeYaml
}

type yaml_codec struct{}

func (y yaml_codec) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yaml_codec) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

type xml_encoder struct{}

func (x xml_encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	b, err := XMLCodec.Marshal(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	w.Set(b)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (x xml_encoder) ContentType() string {
	return MimeXml
}

type xml_codec struct{}

func (x xml_codec) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

func (x xml_codec) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

type proto_encoder struct{}

func (p proto_encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	b, err := PROTOCodec.Marshal(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	w.Set(b)
	r := &closerWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (p proto_encoder) ContentType() string {
	return MimeProto
}

type proto_codec struct{}

func (p proto_codec) Marshal(v any) ([]byte, error) {
	if _, ok := v.(protoreflect.ProtoMessage); !ok {
		return nil, errors.WithStack(errUnsupportedData)
	}
	return proto.Marshal(v.(protoreflect.ProtoMessage))
}

func (p proto_codec) Unmarshal(data []byte, v any) error {
	if _, ok := v.(protoreflect.ProtoMessage); !ok {
		return errors.WithStack(errUnsupportedData)
	}
	return proto.Unmarshal(data, v.(protoreflect.ProtoMessage))
}

type form_encoder struct{}

// Encode do not support float number
func (f form_encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
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
	c.B.Reset()
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
