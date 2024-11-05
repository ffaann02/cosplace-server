package utils

import (
	"errors"
	"reflect"
)

func ValidateStruct(s interface{}) (bool, string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return false, "", errors.New("input is not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			return false, v.Type().Field(i).Name, nil
		}
	}
	return true, "", nil
}
