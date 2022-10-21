package hasaki

import (
	"reflect"
	"strconv"
	"strings"
)

func getKeys(m map[string]interface{}) []string {
	results := make([]string, 0, len(m))
	for k, _ := range m {
		results = append(results, k)
	}
	return results
}

func getValues(m map[string]interface{}) []interface{} {
	results := make([]interface{}, 0, len(m))
	for _, v := range m {
		results = append(results, v)
	}
	return results
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return ""
}

func structToAny(v interface{}, tag string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var types = reflect.TypeOf(v)
	var values = reflect.ValueOf(v)
	switch types.Kind() {
	case reflect.Ptr:
		if values.IsNil() {
			return m, nil
		}
		types = types.Elem()
		values = values.Elem()
		if types.Kind() != reflect.Struct {
			return nil, ErrDataNotSupported
		}
	case reflect.Struct:
	default:
		return nil, ErrDataNotSupported
	}
	doStructToAny(m, types, values, tag)
	return m, nil
}

func doStructToAny(m map[string]interface{}, types reflect.Type, values reflect.Value, tag string) {
	for i := 0; i < values.NumField(); i++ {
		f := values.Field(i)
		t := types.Field(i)
		var key = t.Name
		if !isPublic(key) {
			continue
		}

		switch t.Type.Kind() {
		case reflect.Struct:
			doStructToAny(m, t.Type, f, tag)
			continue
		case reflect.Ptr:
			if f.IsNil() {
				continue
			}
			if t.Type.Elem().Kind() == reflect.Struct {
				doStructToAny(m, t.Type.Elem(), f.Elem(), tag)
				continue
			}
		}

		if tag == "json" {
			s := t.Tag.Get("json")
			if s == "-" {
				continue
			}
			if s != "" {
				if arr := strings.Split(s, ","); len(arr) > 0 {
					key = arr[0]
				}
			}
		} else if tag == "form" {
			s := t.Tag.Get("form")
			if s == "-" {
				continue
			}
			if s != "" {
				key = s
			}
		}
		m[key] = f.Interface()
	}
}

func isPublic(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}
