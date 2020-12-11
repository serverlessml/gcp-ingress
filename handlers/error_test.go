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

package handlers_test

import (
	"testing"

	"github.com/serverlessml/gcp-ingress/handlers"
)

func TestError(t *testing.T) {
	tests := []struct {
		name string
		in   handlers.Error
		want string
	}{
		{
			name: "With details",
			in: handlers.Error{
				Type:    "test",
				Message: "foo",
				Details: map[string]interface{}{"foo": "bar"},
			},
			want: "[test] foo. Details:\nmap[foo:bar]",
		},
		{
			name: "Without details",
			in: handlers.Error{
				Type:    "test",
				Message: "foo",
			},
			want: "[test] foo",
		},
	}

	for _, test := range tests {
		got := test.in.Error()
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
				test.name, test.want, got)
		}
	}
}

func TestErrorPush(t *testing.T) {
	tests := []struct {
		name string
		in   handlers.ErrorPush
		want string
	}{
		{
			name: "With details",
			in: handlers.ErrorPush{
				Message: "foo",
				Details: map[string]interface{}{"foo": "bar"},
			},
			want: "foo. Details:\nmap[foo:bar]",
		},
		{
			name: "Without details",
			in: handlers.ErrorPush{
				Message: "foo",
			},
			want: "foo",
		},
	}

	for _, test := range tests {
		got := test.in.Error()
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
				test.name, test.want, got)
		}
	}
}
