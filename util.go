package hasaki

import (
	"fmt"
	"net/url"
	"strconv"
)

type ContentTypeInterface interface {
	String() string
}

type ContentType string

func (c ContentType) String() string {
	return string(c)
}

const (
	JSON   ContentType = "application/json; charset=utf-8"
	FORM   ContentType = "application/x-www-form-urlencoded"
	STREAM ContentType = "application/octet-stream"
	JPEG   ContentType = "image/jpeg"
	GIF    ContentType = "image/gif"
	PNG    ContentType = "image/png"
	MP4    ContentType = "video/mpeg4"
)

type Any map[string]interface{}

type Form map[string]string

func FormEncode(data Any) string {
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
		} else if v, ok := item.(float64); ok {
			form.Set(k, fmt.Sprintf("%.6f", v))
			continue
		} else if v, ok := item.([]string); ok {
			form[k] = v
			continue
		} else if v, ok := item.([]int); ok {
			for _, num := range v {
				form.Add(k, strconv.Itoa(num))
			}
			continue
		} else if v, ok := item.([]int64); ok {
			for _, num := range v {
				form.Add(k, strconv.Itoa(int(num)))
			}
			continue
		} else if v, ok := item.([]float64); ok {
			for _, num := range v {
				form.Add(k, fmt.Sprintf("%.6f", num))
			}
			continue
		}
	}
	return form.Encode()
}
