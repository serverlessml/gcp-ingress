package validator_test

import (
	"testing"

	"github.com/serverlessml/gcp-ingress/validator"
)

func TestIsSHA1(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{
			name: "Positive",
			in:   "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
			want: true,
		},
		{
			name: "Negative",
			in:   "8c2f3d3c5dd853231c7429b099347d13c8bb2c371",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				got := validator.IsSHA1(test.in)
				if got != test.want {
					t.Fatalf("[%s]: Results don't match\ngot: %v\nwant: %v",
						test.name, got, test.want)
				}
			})
	}
}
