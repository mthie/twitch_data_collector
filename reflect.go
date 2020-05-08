package main

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func fieldsToMap(in interface{}) (map[string]string, error) {
	if reflect.TypeOf(in).Kind() != reflect.Ptr {
		return nil, errors.New("Non-pointer given")
	}

	if kind := reflect.ValueOf(in).Elem().Kind(); kind != reflect.Struct {
		return nil, errors.Errorf("Non-struct given: %s", kind)
	}

	var out = map[string]string{}

	st := reflect.ValueOf(in).Elem()
	for i := 0; i < st.NumField(); i++ {
		valField := st.Field(i)
		typeField := st.Type().Field(i)

		jsonTag := strings.Split(typeField.Tag.Get("json"), ",")[0]
		if jsonTag == "" {
			// Empty tag, skip
			continue
		}

		switch typeField.Type {
		case reflect.TypeOf(time.Time{}):
			out[jsonTag] = valField.Addr().Interface().(*time.Time).String()
			continue
		}

		switch typeField.Type.Kind() {

		case reflect.Bool:
			out[jsonTag] = strconv.FormatBool(valField.Bool())

		case reflect.Int, reflect.Int64:
			out[jsonTag] = strconv.FormatInt(valField.Int(), 10)

		case reflect.String:
			out[jsonTag] = valField.String()

		case reflect.Slice:
			switch typeField.Type.Elem().Kind() {
			case reflect.String:
				res := valField.Addr().Interface().(*[]string)
				out[jsonTag] = strings.Join(*res, "\n")
			}

		default:
			return nil, errors.Errorf("Unhandled field type: %s", typeField.Type.String())

		}
	}

	return out, nil
}
