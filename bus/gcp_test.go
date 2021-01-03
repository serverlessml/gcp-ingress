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

package bus_test

import (
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	bus "github.com/serverlessml/ingress/bus"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var client bus.GCPClient

func getMockServer() *grpc.ClientConn {
	srv := pstest.NewServer()
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	return conn
}

var MockServerOptions []option.ClientOption = []option.ClientOption{option.WithGRPCConn(getMockServer())}

func getClient() (bus.GCPClient, error) {
	c := client
	c.ProjectID = "test"
	c.Opts = MockServerOptions
	err := c.Connect()
	if err != nil {
		return bus.GCPClient{}, err
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
