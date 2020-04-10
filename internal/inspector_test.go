package internal

import (
	"reflect"
	"testing"
)

func TestParseFields(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []ParsedField
		wantErr bool
	}{
		// TODO: more test cases
		{
			name:    "Not a struct value or pointer to struct",
			args:    args{s: []string{"not", "a", "struct"}},
			wantErr: true,
		},
		{
			name: "Struct value with all tags",
			args: args{s: struct {
				Foo bool     `name:"FOO_ENV"`
				Bar int      `name:"BAR_ENV" required:"true"`
				Baz []string `name:"BAZ_ENV" required:"true"`
			}{}},
			want: []ParsedField{
				{EnvName: "FOO_ENV", EnvRequired: false, FieldName: "Foo", FieldType: reflect.TypeOf(false)},
				{EnvName: "BAR_ENV", EnvRequired: true, FieldName: "Bar", FieldType: reflect.TypeOf(0)},
				{EnvName: "BAZ_ENV", EnvRequired: true, FieldName: "Baz", FieldType: reflect.TypeOf([]string{""}), ElemType: reflect.TypeOf("")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFields(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testStruct struct {
	Foo bool
	Bar int
	Baz []string
}

func TestSetValues(t *testing.T) {
	type args struct {
		s      testStruct
		values map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    testStruct
		wantErr bool
	}{
		{
			name: "Valid values",
			args: args{s: testStruct{}, values: map[string]interface{}{"Foo": true, "Bar": 5}},
			want: testStruct{
				Foo: true,
				Bar: 5,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := SetValues(&test.args.s, test.args.values)
			if (err != nil) != test.wantErr {
				t.Errorf("SetValues() error = %v, wantErr %v", err, test.wantErr)
			}
			if !reflect.DeepEqual(test.args.s, test.want) {
				t.Errorf("SetValues() = %v, want %v", test.args.s, test.want)
			}
		})
	}
}
