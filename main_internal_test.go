package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/serverlessml/gcp-ingress/processor"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/serverlessml/gcp-ingress/bus"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func TestGetEnv(t *testing.T) {
	type args struct {
		key, fallback string
	}

	os.Setenv("TEST_TestGetEnv", "test")

	tests := []struct {
		name string
		in   *args
		want string
	}{
		{
			name: "From env",
			in:   &args{key: "TEST_TestGetEnv", fallback: ""},
			want: "test",
		},
		{
			name: "From fallback",
			in:   &args{key: "TEST_TestGetEnv1", fallback: "test"},
			want: "test",
		},
	}
	for _, test := range tests {
		got := GetEnv(test.in.key, test.in.fallback)
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
				test.name, test.want, got)
		}
	}
	os.Unsetenv("TEST_TestGetEnv")
}

func TestGetRequestPayload(t *testing.T) {
	type output struct {
		data []byte
		err  error
	}

	tests := []struct {
		name string
		in   io.ReadCloser
		want *output
	}{
		{
			name: "Positive",
			in:   ioutil.NopCloser(strings.NewReader("test")),
			want: &output{
				data: []byte("test"),
				err:  nil,
			},
		},
	}
	for _, test := range tests {
		got := GetRequestPayload(test.in)
		if !reflect.DeepEqual(got, test.want.data) {
			t.Fatalf("[%s]: Results don't match\nwant: %v\ngot: %v",
				test.name, test.want.data, got)
		}
	}
}

func TestGetMustMarshal(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want []byte
	}{
		{
			name: "Positive",
			in:   `{"foo": "bar"}`,
			want: []byte{34, 123, 92, 34, 102, 111, 111, 92, 34, 58, 32, 92, 34, 98, 97, 114, 92, 34, 125, 34},
		},
	}
	for _, test := range tests {
		got := MustMarshal(test.in)

		if test.name == "Positive" {
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("[%s]: Results don't match\nwant: %v\ngot: %v",
					test.name, test.want, got)
			}
		} else {
			if got == nil {
				t.Fatalf("[%s]: Wrong error implementation", test.name)
			}
		}
	}
}

func getMockServer() *grpc.ClientConn {
	srv := pstest.NewServer()
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	return conn
}

func getClient() bus.Client {
	var c bus.Client
	c.ProjectID = "test"
	c.Connect(option.WithGRPCConn(getMockServer()))
	return c
}

func TestRunner(t *testing.T) {
	pubsubClient = getClient()

	probe := []byte{123, 10, 32, 32, 32, 32, 34, 112, 114, 111, 106, 101, 99, 116, 95, 105, 100, 34, 58, 32, 34, 48, 99, 98, 97, 56, 50, 102, 102, 45, 57, 55, 57, 48, 45, 52, 53, 52, 100, 45, 98, 55, 98, 57, 45, 50, 50, 53, 55, 48, 101, 55, 98, 97, 50, 56, 99, 34, 44, 10, 32, 32, 32, 32, 34, 99, 111, 100, 101, 95, 104, 97, 115, 104, 34, 58, 32, 34, 56, 99, 50, 102, 51, 100, 51, 99, 53, 100, 100, 56, 53, 51, 50, 51, 49, 99, 55, 52, 50, 57, 98, 48, 57, 57, 51, 52, 55, 100, 49, 51, 99, 56, 98, 98, 50, 99, 51, 55, 34, 44, 10, 32, 32, 32, 32, 34, 112, 105, 112, 101, 108, 105, 110, 101, 95, 99, 111, 110, 102, 105, 103, 34, 58, 32, 91, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 100, 97, 116, 97, 34, 58, 32, 123, 34, 102, 111, 111, 34, 58, 32, 49, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 109, 111, 100, 101, 108, 34, 58, 32, 123, 34, 102, 111, 111, 34, 58, 32, 49, 125, 10, 32, 32, 32, 32, 32, 32, 32, 32, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 100, 97, 116, 97, 34, 58, 32, 123, 34, 102, 111, 111, 34, 58, 32, 50, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 109, 111, 100, 101, 108, 34, 58, 32, 123, 34, 98, 97, 114, 34, 58, 32, 50, 125, 10, 32, 32, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 93, 10, 125}

	tests := []struct {
		name    string
		in      []byte
		want    OutputPayload
		isError bool
	}{
		{
			name: "Positive: submission errors returned.",
			in:   probe,
			want: OutputPayload{
				Errors: []errorOutput{
					{
						PipelineConfig: processor.PipelineConfig{Data: map[string]interface{}{"foo": 1}, Model: map[string]interface{}{"foo": 1}},
						Message:        `rpc error: code = NotFound desc = topic "projects/test/topics/0cba82ff-9790-454d-b7b9-22570e7ba28c"`,
					},
					{
						PipelineConfig: processor.PipelineConfig{Data: map[string]interface{}{"foo": 2}, Model: map[string]interface{}{"bar": 2}},
						Message:        `rpc error: code = NotFound desc = topic "projects/test/topics/0cba82ff-9790-454d-b7b9-22570e7ba28c"`,
					},
				},
				SubmittedID: []string{},
			},
			isError: false,
		},
		{
			name: "Positive",
			in:   probe,
			want: OutputPayload{
				Errors: []errorOutput{},
				SubmittedID: []string{
					"322ededf-4587-4c08-a5ee-a177308601ef",
					"beca0bb7-aafa-4d30-b528-d7a6b5694c23",
				},
			},
			isError: false,
		},
		{
			name:    "Negative: proc.Exec",
			in:      []byte{1},
			want:    OutputPayload{},
			isError: true,
		},
	}

	for _, test := range tests {
		if test.name == "Positive" && !test.isError {
			// create fake topic with the name of the project
			pubsubClient.Instance.CreateTopic(pubsubClient.Ctx, "0cba82ff-9790-454d-b7b9-22570e7ba28c")
		}
		got, err := runner(test.in)

		if !test.isError {
			if err != nil {
				t.Fatalf("[%s]: Error: %s", test.name, err)
			}
			if test.name == "Positive" {
				// hardcode UUID
				got.SubmittedID = []string{
					"322ededf-4587-4c08-a5ee-a177308601ef",
					"beca0bb7-aafa-4d30-b528-d7a6b5694c23",
				}
			}
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", test.want) {
				t.Fatalf("[%s]: Results don't match\nwant: %v\ngot: %v",
					test.name, test.want, got)
			}
		} else {
			if err == nil {
				t.Fatalf("[%s]: Wrong error implementation", test.name)
			}
		}
	}
}
