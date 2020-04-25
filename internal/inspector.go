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
	FieldIndex  []int
	FieldType   reflect.Type
	ElemType    reflect.Type
	EnvName     string
	EnvRequired bool
}

func ParseFields(s interface{}, nestedIndex []int) ([]ParsedField, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, ErrIsNotStruct
	}

	n := v.NumField()
	fields := make([]ParsedField, 0, n)
	for i := 0; i < n; i++ {
		if v.Field(i).Kind() == reflect.Struct {
			ff, err := ParseFields(v.Field(i).Interface(), append(nestedIndex, i))
			if err != nil {
				return nil, err
			}
			fields = append(fields, ff...)
			continue
		}

		structField := v.Type().Field(i)

		var field ParsedField
		field.FieldType = structField.Type
		if field.FieldType.Kind() == reflect.Slice {
			field.ElemType = structField.Type.Elem()
		}
		field.FieldIndex = append(nestedIndex, i)

		field.EnvName, _ = structField.Tag.Lookup(envName)
		required, _ := structField.Tag.Lookup(envRequired)
		field.EnvRequired, _ = strconv.ParseBool(required)

		fields = append(fields, field)
	}

	return fields, nil
}

func SetValues(s interface{}, values []FieldValue) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Ptr {
		return ErrIsNotPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	v := reflect.ValueOf(s).Elem()
	for _, value := range values {
		if reflect.ValueOf(value.Value).IsValid() {
			v.FieldByIndex(value.StructIndex).Set(reflect.ValueOf(value.Value))
		}
	}

	return nil
}
