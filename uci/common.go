package uci

import (
	"fmt"
	"reflect"
)

// StructToMap converts a struct to a map with json tag as key
func StructToMap(s interface{}) (map[string]string, error) {
	m := map[string]string{}
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("json") == "" {
			continue
		}
		// if value is "" then skip
		if fmt.Sprintf("%v", v.Field(i).Interface()) == "" {
			continue
		}
		m[f.Tag.Get("json")] = fmt.Sprintf("%v", v.Field(i).Interface())
	}
	return m, nil
}
