package hasaki

import (
	"net/url"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type Any map[string]interface{}

type H map[string]string

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
	var data Any
	switch val := v.(type) {
	case Any:
		data = val
	case map[string]interface{}:
		data = val
	default:
		if formData, err := structToAny(v, "form"); err != nil {
			return nil, err
		} else {
			data = formData
		}
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
