package hasaki

import (
	"bytes"
	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"strings"
)

const (
	ContentTypeJSON   = "application/json;charset=utf-8"
	ContentTypeFORM   = "application/x-www-form-urlencoded"
	ContentTypeSTREAM = "application/octet-stream"
	ContentTypeJPEG   = "image/jpeg"
	ContentTypeGIF    = "image/gif"
	ContentTypePNG    = "image/png"
	ContentTypeMP4    = "video/mpeg4"
)

type Any map[string]interface{}

type H map[string]string

type Encoder interface {
	Encode(v interface{}) (io.Reader, error)
	GetContentType() string
}

var (
	JsonEncoder = new(json_encoder)
	FormEncoder = new(form_encoder)
)

type json_encoder struct{}

func (j json_encoder) Encode(v interface{}) (io.Reader, error) {
	w := &bytes.Buffer{}
	err := jsoniter.NewEncoder(w).Encode(v)
	return w, errors.WithStack(err)
}

func (j json_encoder) GetContentType() string {
	return ContentTypeJSON
}

type form_encoder struct{}

// Encode do not support float number
func (f form_encoder) Encode(v interface{}) (io.Reader, error) {
	values, err := form.NewEncoder().Encode(v)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return strings.NewReader(values.Encode()), nil
}

func (f form_encoder) GetContentType() string {
	return ContentTypeFORM
}
