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

package handlers_test

import (
	"testing"

	bus "github.com/serverlessml/ingress/handlers"
)

type inputPush struct {
	Payload []byte
	Topic   string
	Ch      chan error
}

func TestPushRoutine(t *testing.T) {
	c := getClient("test", "foo")

	tests := []struct {
		name    string
		in      *inputPush
		want    error
		isError bool
	}{
		{
			name: "Positive",
			in: &inputPush{
				Payload: []byte{1},
				Topic:   "foo",
				Ch:      make(chan error, 1),
			},
			want:    nil,
			isError: false,
		},
	}

	for _, test := range tests {
		bus.PushRoutine(c, test.in.Payload, test.in.Topic, test.in.Ch)
		got := <-test.in.Ch
		if !test.isError {
			if got != test.want {
				t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
					test.name, test.want, got)
			}
		} else {
			if got == nil {
				t.Fatalf("[%s]: Error wasn't thrown.", test.name)
			}
		}
	}
}
