package internal

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadEnvironment(t *testing.T) {
	type args struct {
		fields []ParsedField
	}
	tests := []struct {
		name    string
		envs    map[string]string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "With all vars set and required",
			envs: map[string]string{"STRING_ENV": "foo", "INT_ENV": "-5", "UINT_ENV": "5", "FLOAT_ENV": "1.23",
				"BOOL_ENV": "true", "SLICE_ENV": "a,b,c", "DURATION_ENV": "1h15m11s"},
			args: args{
				fields: []ParsedField{
					{EnvName: "STRING_ENV", EnvRequired: true, FieldName: "String", FieldType: reflect.TypeOf("")},
					{EnvName: "INT_ENV", EnvRequired: true, FieldName: "Int", FieldType: reflect.TypeOf(0)},
					{EnvName: "UINT_ENV", EnvRequired: true, FieldName: "Uint", FieldType: reflect.TypeOf(uint(0))},
					{EnvName: "FLOAT_ENV", EnvRequired: true, FieldName: "Float", FieldType: reflect.TypeOf(0.00)},
					{EnvName: "BOOL_ENV", EnvRequired: true, FieldName: "Bool", FieldType: reflect.TypeOf(true)},
					{EnvName: "SLICE_ENV", EnvRequired: true, FieldName: "Slice", FieldType: reflect.TypeOf([]string{""}), ElemType: reflect.TypeOf("")},
					{EnvName: "DURATION_ENV", EnvRequired: true, FieldName: "Duration", FieldType: reflect.TypeOf(time.Duration(0))},
				},
			},
			want: map[string]interface{}{"String": "foo", "Int": -5, "Uint": uint(5), "Float": 1.23,
				"Bool": true, "Slice": []string{"a", "b", "c"}, "Duration": time.Hour + 15*time.Minute + 11*time.Second},
		},
		{
			name: "With unset not required vars",
			envs: map[string]string{"FOO_ENV": "foo"},
			args: args{
				fields: []ParsedField{
					{EnvName: "FOO_ENV", EnvRequired: true, FieldName: "Foo", FieldType: reflect.TypeOf("")},
					{EnvName: "BAR_ENV", EnvRequired: false, FieldName: "Bar", FieldType: reflect.TypeOf("")},
				},
			},
			want: map[string]interface{}{"Foo": "foo", "Bar": ""},
		},
		{
			name: "With unset required vars",
			envs: map[string]string{"FOO_ENV": "foo"},
			args: args{
				fields: []ParsedField{
					{EnvName: "FOO_ENV", EnvRequired: true, FieldName: "Foo", FieldType: reflect.TypeOf("")},
					{EnvName: "BAR_ENV", EnvRequired: true, FieldName: "Bar", FieldType: reflect.TypeOf("")},
				},
			},
			wantErr: true,
		},
		{
			name: "Wrong type in var",
			envs: map[string]string{"FOO_ENV": "foo", "BAR_ENV": "bar"},
			args: args{
				fields: []ParsedField{
					{EnvName: "FOO_ENV", EnvRequired: true, FieldName: "Foo", FieldType: reflect.TypeOf("")},
					{EnvName: "BAR_ENV", EnvRequired: true, FieldName: "Bar", FieldType: reflect.TypeOf(5)},
				},
			},
			wantErr: true,
		},
		{
			name: "Unknown type in var",
			envs: map[string]string{"FOO_ENV": "foo", "BAR_ENV": "bar:1"},
			args: args{
				fields: []ParsedField{
					{EnvName: "FOO_ENV", EnvRequired: true, FieldName: "Foo", FieldType: reflect.TypeOf("")},
					{EnvName: "BAR_ENV", EnvRequired: true, FieldName: "Bar", FieldType: reflect.TypeOf(map[string]int{})},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanup := setEnvironment(t, test.envs)
			defer cleanup()

			got, err := ReadEnvironment(test.args.fields)
			if (err != nil) != test.wantErr {
				t.Errorf("ReadEnvironment() error = %v, wantErr %v", errors.Unwrap(err), test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ReadEnvironment() = %v, want %v", got, test.want)
			}
		})
	}
}

func setEnvironment(t *testing.T, values map[string]string) (cleanup func()) {
	for name, value := range values {
		err := os.Setenv(name, value)
		if err != nil {
			t.Logf("failed to set %q env: %v", name, err)
		}
	}

	return func() {
		for name := range values {
			err := os.Unsetenv(name)
			if err != nil {
				t.Logf("failed to unset %q env: %v", name, err)
			}
		}
	}
}
