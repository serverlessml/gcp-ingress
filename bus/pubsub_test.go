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

func getClient() bus.Client {
	c := client
	c.ProjectID = "test"
	c.Connect(option.WithGRPCConn(getMockServer()))
	c.Instance.CreateTopic(c.Ctx, "foo")
	return c
}

type input struct {
	Payload []byte
	Topic   string
	Ch      chan error
}

func TestConnect(t *testing.T) {
	mockServer := getMockServer()
	defer mockServer.Close()

	c := client
	c.ProjectID = "test"
	err := c.Connect(option.WithGRPCConn(mockServer))
	defer c.Instance.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPush(t *testing.T) {
	c := getClient()
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
	c := getClient()
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
