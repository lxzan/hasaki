package hasaki

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type Any map[string]interface{}

type Form map[string]string

func FormEncode(data Any) (string, error) {
	var form = url.Values{}
	for k, item := range data {
		switch item.(type) {
		case string:
			form.Set(k, item.(string))
		case int:
			form.Set(k, strconv.Itoa(item.(int)))
		case int64:
			form.Set(k, strconv.Itoa(int(item.(int64))))
		case float64:
			form.Set(k, fmt.Sprintf("%.9f", item.(float64)))
		case []string:
			var arr = item.([]string)
			for _, str := range arr {
				form.Add(k, str)
			}
		case []int:
			var arr = item.([]int)
			for _, num := range arr {
				form.Add(k, strconv.Itoa(num))
			}
		case []int64:
			var arr = item.([]int64)
			for _, num := range arr {
				form.Add(k, strconv.Itoa(int(num)))
			}
		case []float64:
			var arr = item.([]float64)
			for _, num := range arr {
				form.Add(k, fmt.Sprintf("%.9f", num))
			}
		default:
			return "", errors.New("unsupported data type")
		}
	}
	return form.Encode(), nil
}
