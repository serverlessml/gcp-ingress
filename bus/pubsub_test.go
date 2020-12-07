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

package bus_test

import (
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/serverlessml/gcp-ingress/bus"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var client bus.Client

func getMockServer() *grpc.ClientConn {
	srv := pstest.NewServer()
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	return conn
}

var MockServerOptions []option.ClientOption = []option.ClientOption{option.WithGRPCConn(getMockServer())}

func getClient() (bus.Client, error) {
	c := client
	c.ProjectID = "test"
	c.Opts = MockServerOptions
	err := c.Connect()
	if err != nil {
		return bus.Client{}, err
	}
	c.Instance.CreateTopic(c.Ctx, "foo")
	return c, nil
}

type input struct {
	Payload []byte
	Topic   string
	Ch      chan error
}

func TestConnect(t *testing.T) {
	mockServer := getMockServer()
	defer mockServer.Close()

	_, err := getClient()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPush(t *testing.T) {
	c, _ := getClient()

	tests := []struct {
		name    string
		in      *input
		want    error
		isError bool
	}{
		{
			name: "Positive",
			in: &input{
				Payload: []byte{1},
				Topic:   "foo",
			},
			want:    nil,
			isError: false,
		},
		{
			name: "Negative: no topic",
			in: &input{
				Payload: []byte{1},
				Topic:   "bar",
			},
			want:    nil,
			isError: true,
		},
	}

	for _, test := range tests {
		err := c.Push(test.in.Payload, test.in.Topic)
		if !test.isError {
			if err != test.want {
				t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
					test.name, test.want, err)
			}
		} else {
			if err == nil {
				t.Fatalf("[%s]: Error wasn't thrown.", test.name)
			}
		}
	}
}

func TestPushRoutine(t *testing.T) {
	c, _ := getClient()
	defer c.Instance.Close()

	tests := []struct {
		name    string
		in      *input
		want    error
		isError bool
	}{
		{
			name: "Positive",
			in: &input{
				Payload: []byte{1},
				Topic:   "foo",
				Ch:      make(chan error, 1),
			},
			want:    nil,
			isError: false,
		},
	}

	for _, test := range tests {
		c.PushRoutine(test.in.Payload, test.in.Topic, test.in.Ch)
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
