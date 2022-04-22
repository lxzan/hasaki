package hasaki

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"net/url"
	"reflect"
)

type Any map[string]interface{}

type Form map[string]string

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	GetContentType() string
}

var (
	JsonEncoder = new(json_encoder)
	FormEncoder = new(form_encoder)
)

type json_encoder struct{}

func (j json_encoder) Encode(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (j json_encoder) GetContentType() string {
	return ContentType_JSON.String()
}

type form_encoder struct{}

// Encode do not support float number
func (f form_encoder) Encode(v interface{}) ([]byte, error) {
	data, ok := v.(Any)
	if !ok {
		return nil, errors.New("only support hasaki.Any type")
	}
	if len(data) == 0 {
		return []byte(""), nil
	}

	var form = url.Values{}
	for k, item := range data {
		if val := reflect.ValueOf(item); val.Kind() == reflect.Slice {
			var length = val.Len()
			for i := 0; i < length; i++ {
				form.Add(k, ToString(val.Index(i).Interface()))
			}
		} else {
			form.Set(k, ToString(item))
		}
	}

	return []byte(form.Encode()), nil
}

func (f form_encoder) GetContentType() string {
	return ContentType_FORM.String()
}
