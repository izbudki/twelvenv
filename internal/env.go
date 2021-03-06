package internal

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	duration = reflect.Kind(100)
)

var (
	ErrNotSet        = errors.New("variable not set but required")
	ErrCanNotConvert = errors.New("can not convert")
)

type FieldValue struct {
	StructIndex []int
	Value       interface{}
}

func ReadEnvironment(fields []ParsedField) ([]FieldValue, error) {
	values := make([]FieldValue, len(fields))
	for i, field := range fields {
		value, ok := os.LookupEnv(field.EnvName)
		if !ok && field.EnvRequired {
			return nil, fmt.Errorf("lookup %s: %w", field.EnvName, ErrNotSet)
		}
		converted, err := convertToBasicType(value, field.FieldType, field.ElemType)
		if err != nil {
			return nil, fmt.Errorf("convert %q variable value: %w", field.EnvName, err)
		}
		values[i] = FieldValue{StructIndex: field.FieldIndex, Value: converted}
	}
	return values, nil
}

func convertToBasicType(value string, t reflect.Type, e reflect.Type) (interface{}, error) {
	kind := t.Kind()
	if t.PkgPath() == "time" && t.Name() == "Duration" {
		kind = duration
	}

	switch kind {
	case reflect.String:
		return value, nil
	case reflect.Int:
		converted, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("string %s to %s: %w", value, t, err)
		}
		return int(converted), nil
	case reflect.Uint:
		converted, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("string %s to %s: %w", value, t, err)
		}
		return uint(converted), nil
	case reflect.Float64:
		converted, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("string %s to %s: %w", value, t, err)
		}
		return converted, nil
	case reflect.Bool:
		converted, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("string %s to %s: %w", value, t, err)
		}
		return converted, nil
	case reflect.Slice:
		values := strings.Split(value, ",")
		if len(values) == 1 && values[0] == "" {
			return nil, nil
		}

		converted := reflect.MakeSlice(reflect.SliceOf(e), 0, len(values))
		for i := range values {
			elem, err := convertToBasicType(values[i], e, nil)
			if err != nil {
				return nil, err
			}
			converted = reflect.Append(converted, reflect.ValueOf(elem))
		}
		return converted.Interface(), nil
	case duration:
		converted, err := time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("string %s to %s: %w", value, t, err)
		}
		return converted, nil
	default:
		return nil, fmt.Errorf("string %q to %s: %w", value, t, ErrCanNotConvert)
	}
}
