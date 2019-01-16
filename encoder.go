package hasaki

import (
	"bytes"
	"github.com/json-iterator/go"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type ContentType string

type Encoder interface {
	GetReader() io.Reader
	GetContentType() string
}

type JsonEncoder struct {
	reader      io.Reader
	contentType string
}

func (this *JsonEncoder) GetReader() io.Reader {
	return this.reader
}

func (this *JsonEncoder) GetContentType() string {
	return "application/json; charset=utf-8"
}

func NewJsonEncoder(data JSON) *JsonEncoder {
	b, _ := jsoniter.Marshal(data)
	return &JsonEncoder{
		reader: bytes.NewReader(b),
	}
}

type FormEncoder struct {
	reader      io.Reader
	contentType string
}

func (this *FormEncoder) GetReader() io.Reader {
	return this.reader
}

func (this *FormEncoder) GetContentType() string {
	return "application/x-www-form-urlencoded"
}

func NewFormEncoder(data JSON) *FormEncoder {
	return &FormEncoder{
		reader: strings.NewReader(querystring(data)),
	}
}

func querystring(data JSON) string {
	var form = url.Values{}
	for k, item := range data {
		if v, ok := item.(string); ok {
			form.Set(k, v)
			continue
		} else if v, ok := item.(int); ok {
			form.Set(k, strconv.Itoa(v))
			continue
		} else if v, ok := item.(int64); ok {
			form.Set(k, strconv.Itoa(int(v)))
			continue
		} else if v, ok := item.([]string); ok {
			form[k+"[]"] = v
			continue
		} else if v, ok := item.([]int); ok {
			var key = k + "[]"
			for _, num := range v {
				form.Add(key, strconv.Itoa(num))
			}
			continue
		} else if v, ok := item.([]int64); ok {
			var key = k + "[]"
			for _, num := range v {
				form.Add(key, strconv.Itoa(int(num)))
			}
			continue
		}
	}
	return form.Encode()
}
