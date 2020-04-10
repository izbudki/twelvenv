package internal

import (
	"errors"
	"reflect"
	"strconv"
)

const (
	envName     = "name"
	envRequired = "required"
)

var (
	ErrIsNotStruct  = errors.New("is not struct")
	ErrIsNotPointer = errors.New("is not pointer")
)

type ParsedField struct {
	FieldName   string
	FieldType   reflect.Type
	ElemType    reflect.Type
	EnvName     string
	EnvRequired bool
}

func ParseFields(s interface{}) ([]ParsedField, error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, ErrIsNotStruct
	}

	n := t.NumField()
	fields := make([]ParsedField, n)
	for i := 0; i < n; i++ {
		fields[i].FieldName = t.Field(i).Name
		fields[i].FieldType = t.Field(i).Type
		if fields[i].FieldType.Kind() == reflect.Slice {
			fields[i].ElemType = t.Field(i).Type.Elem()
		}

		fields[i].EnvName, _ = t.Field(i).Tag.Lookup(envName)
		required, _ := t.Field(i).Tag.Lookup(envRequired)
		fields[i].EnvRequired, _ = strconv.ParseBool(required)
	}

	return fields, nil
}

func SetValues(s interface{}, values map[string]interface{}) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return ErrIsNotPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	v := reflect.ValueOf(s).Elem()
	for name, value := range values {
		if reflect.ValueOf(value).IsValid() {
			v.FieldByName(name).Set(reflect.ValueOf(value))
		}
	}

	return nil
}
