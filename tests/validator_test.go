package tests

import (
	"reflect"
	"testing"

	"github.com/serverlessml/gcp-ingress/validator"
)

func TestValidator(t *testing.T) {
	type testStruct struct {
		Foo string `validate:"required,sha1"`
	}

	tests := []struct {
		name string
		in   testStruct
		want string
	}{
		{
			name: "Positive",
			in:   testStruct{Foo: "8c2f3d3c5dd853231c7429b099347d13c8bb2c37"},
			want: "",
		},
		{
			name: "Negative",
			in:   testStruct{Foo: "8c2f3d3c5dd853231c7429b099347d13c8bb2c371"},
			want: `Key: 'testStruct.Foo' Error:Field validation for 'Foo' failed on the 'sha1' tag`,
		},
	}

	Validator := validator.New()

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				got := Validator.Struct(test.in)
				if reflect.DeepEqual(validator.GetValidationErrors(got), test.want) {
					t.Fatalf("[%s]: Results don't match\ngot: %s\nwant: %s",
						test.name, got, test.want)
				}
			})
	}
}
