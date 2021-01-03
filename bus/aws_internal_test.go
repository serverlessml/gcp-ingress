/*
Copyright © 2020 Dmitry Kisler <admin@dkisler.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bus

import (
	"testing"
)

func TestTopic(t *testing.T) {
	type inpt struct {
		Region  string
		Project string
		Topic   string
	}
	tests := []struct {
		in   *inpt
		want string
	}{
		{
			in: &inpt{
				Region:  "eu-west-1",
				Project: "111111",
				Topic:   "foo",
			},
			want: "arn:aws:sns:eu-west-1:111111:foo",
		},
	}

	for _, test := range tests {
		c := &AWSClient{Region: test.in.Region, ProjectID: test.in.Project}
		got := c.getTopicArn(test.in.Topic)
		if got != test.want {
			t.Fatalf("Results don't match\nwant: %s\ngot: %s", test.want, test.in)
		}
	}
}
