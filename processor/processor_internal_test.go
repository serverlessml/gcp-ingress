package processor

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name      string
		in        []byte
		want      Input
		wantError error
	}{
		{
			name: "Positive",
			in: []byte(`{
					"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c",
					"code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
					"pipeline_config": {"data": {}, "mode": {}},
				}`),
			want: Input{
				ProjectID: "0cba82ff-9790-454d-b7b9-22570e7ba28c",
				CodeHash:  "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
				Config: []PipelineConfig{
					{
						Data:  map[string]interface{}{},
						Model: map[string]interface{}{},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Negative",
			in: []byte(`{
					"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c",
					"code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
					"pipeline_config": {"data": {}, "mode": {},
				}`),
			want:      Input{},
			wantError: fmt.Errorf("invalid character '}' looking for beginning of object key string"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				got, err := readInput(test.in)
				if test.name == "Positive" {
					if !reflect.DeepEqual(got, test.want) {
						t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
							test.name, test.want, got)
					}
				} else {
					if err.Error() != test.wantError.Error() {
						t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
							test.name, test.wantError, err)
					}
				}
			})
	}
}
