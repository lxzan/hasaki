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

const (
	JsonContent ContentType = "application/json; charset=utf-8"
	FormContent ContentType = "application/x-www-form-urlencoded"
)

func JsonEncode(data JSON) (ContentType, io.Reader) {
	b, _ := jsoniter.Marshal(data)
	return JsonContent, bytes.NewReader(b)
}

func FormEncode(data JSON) (ContentType, io.Reader) {
	return FormContent, strings.NewReader(querystring(data))
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
