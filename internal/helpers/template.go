package helpers

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/kuma-framework/kuma/v2/pkg/functions"
)

func convertValue(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Ptr:
		// If the pointer is nil, return nil
		if val.IsNil() {
			return nil
		}
		return convertValue(val.Elem().Interface())
	case reflect.Map:
		m := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			m[keyStr] = convertValue(val.MapIndex(key).Interface())
		}
		return m
	case reflect.String:
		str := val.String()
		return str
	default:
		// For slices or arrays, process each element
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			length := val.Len()
			slice := make([]interface{}, length)
			for i := 0; i < length; i++ {
				slice[i] = convertValue(val.Index(i).Interface())
			}
			return slice
		}
		return value
	}
}

func ReplaceVars(text string, vars interface{}, funcs template.FuncMap) (string, error) {
	vars = convertValue(vars)
	t, err := template.New("").Funcs(functions.GetFuncMap()).Parse(text)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	err = t.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
