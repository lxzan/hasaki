package hasaki

import (
	"net/url"
	"strconv"
)

const (
	JsonType = "application/json; charset=utf-8"
	FormType = "application/x-www-form-urlencoded"
)

type JSON map[string]interface{}

type Form map[string]string

func FormEncode(data JSON) string {
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
