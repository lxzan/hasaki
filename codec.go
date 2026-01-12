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
	// MimeJson JSON MIME类型
	// JSON MIME type
	MimeJson = "application/json;charset=utf-8"
	// MimeYaml YAML MIME类型
	// YAML MIME type
	MimeYaml = "application/x-yaml;charset=utf-8"
	// MimeXml XML MIME类型
	// XML MIME type
	MimeXml = "application/xml;charset=utf-8"
	// MimeProtoBuf Protobuf MIME类型
	// Protobuf MIME type
	MimeProtoBuf = "application/x-protobuf"
	// MimeForm 表单 MIME类型
	// Form MIME type
	MimeForm = "application/x-www-form-urlencoded"
	// MimeStream 流式 MIME类型
	// Stream MIME type
	MimeStream = "application/octet-stream"
	// MimeJpeg JPEG图片 MIME类型
	// JPEG image MIME type
	MimeJpeg = "image/jpeg"
	// MimeGif GIF图片 MIME类型
	// GIF image MIME type
	MimeGif = "image/gif"
	// MimePng PNG图片 MIME类型
	// PNG image MIME type
	MimePng = "image/png"
	// MimeMp4 MP4视频 MIME类型
	// MP4 video MIME type
	MimeMp4 = "video/mpeg4"
)

// Any 通用类型映射
// Generic type map
type Any map[string]any

type (
	// Codec 编解码器接口，同时包含编码和解码功能
	// Codec interface, includes both encoding and decoding functionality
	Codec interface {
		Encoder
		Decoder
	}

	// Encoder 编码器接口
	// Encoder interface
	Encoder interface {
		Encode(v any) (io.Reader, error)
		ContentType() string
	}

	// Decoder 解码器接口
	// Decoder interface
	Decoder interface {
		Decode(r io.Reader, v any) error
	}
)

var (
	// JsonCodec JSON编解码器
	// JSON codec
	JsonCodec = new(jsonCodec)
	// FormCodec 表单编解码器
	// Form codec
	FormCodec = new(formCodec)
	// XmlCodec XML编解码器
	// XML codec
	XmlCodec = new(xmlCodec)
)

type (
	jsonCodec struct{}
	formCodec struct{}
	xmlCodec  struct{}
)

func (c jsonCodec) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	err := jsoniter.ConfigFastest.NewEncoder(w).Encode(v)
	r := &internal.CloserWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (c jsonCodec) ContentType() string {
	return MimeJson
}

func (c jsonCodec) Decode(r io.Reader, v any) error {
	return jsoniter.ConfigFastest.NewDecoder(r).Decode(v)
}

func (f formCodec) Encode(v any) (io.Reader, error) {
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

func (f formCodec) ContentType() string {
	return MimeForm
}

func (f formCodec) Decode(r io.Reader, v any) error {
	values, ok := v.(*url.Values)
	if !ok {
		return errors.Wrap(errUnsupportedData, "v must be *url.Values type")
	}
	var builder = &strings.Builder{}
	var temp = internal.GetBuffer()
	_, _ = io.CopyBuffer(builder, r, temp.Bytes()[:internal.BufferSize])
	internal.PutBuffer(temp)
	result, err := url.ParseQuery(builder.String())
	if err != nil {
		return errors.WithStack(err)
	}
	*values = result
	return nil
}

func (c xmlCodec) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	err := xml.NewEncoder(w).Encode(v)
	r := &internal.CloserWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (c xmlCodec) ContentType() string {
	return MimeXml
}

func (c xmlCodec) Decode(r io.Reader, v any) error {
	return errors.WithStack(xml.NewDecoder(r).Decode(v))
}

type streamEncoder struct {
	contentType string
}

// NewStreamEncoder 创建流式编码器
// Create a stream encoder
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
