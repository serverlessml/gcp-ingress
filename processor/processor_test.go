// Copyright 2020 dkisler.com Dmitry Kisler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND
// NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS
// BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF, OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package processor_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/serverlessml/gcp-ingress/processor"
)

func mustMarshal(obj interface{}) []byte {
	out, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return out
}

func TestExec(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    processor.Output
		isError bool
	}{
		{
			name: "Positive",
			in:   []byte(`{"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c", "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37", "pipeline_config": [{"data": {}, "model": {}}]}`),
			want: processor.Output{
				Payload: []processor.OutputPayload{{
					CodeHash: "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
					RunID:    "0cba82ff-9790-454d-b7b9-22570e7ba28c",
					Config: processor.PipelineConfig{
						Data:  map[string]interface{}{},
						Model: map[string]interface{}{},
					},
				}},
				Distribution: processor.OutputDistribution{
					Topic: "trigger_0cba82ff-9790-454d-b7b9-22570e7ba28c",
				},
			},
			isError: false,
		},
		{
			name:    "Negative: json parsing error",
			in:      []byte(`{"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c"`),
			want:    processor.Output{},
			isError: true,
		},
		{
			name:    "Negative: validation error",
			in:      []byte(`{"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c", "code_hash": "foobar", "pipeline_config": [{"data": {}, "model": {}}]}`),
			want:    processor.Output{},
			isError: true,
		},
	}

	var proc processor.Processor
	proc.TopicPrefix = "trigger_"

	for _, test := range tests {
		got, err := proc.Exec(test.in)
		if test.name == "Positive" {
			if err != nil {
				t.Fatalf("[%s]: Error: %s", test.name, err)
			}
			got.Payload[0].RunID = "0cba82ff-9790-454d-b7b9-22570e7ba28c"
			if !reflect.DeepEqual(*got, test.want) {
				t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
					test.name, mustMarshal(test.want), mustMarshal(got))
			}
		} else {
			if err == nil {
				t.Fatalf("[%s]: Error wasn't thrown.", test.name)
			}
		}
	}
}
