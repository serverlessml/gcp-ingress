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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/serverlessml/gcp-ingress/bus"
	"github.com/serverlessml/gcp-ingress/config"
	"github.com/serverlessml/gcp-ingress/handlers"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func mustMarshal(obj interface{}) []byte {
	out, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return out
}

var client bus.GCPClient
var MockServerOptions []option.ClientOption = []option.ClientOption{option.WithGRPCConn(getMockServer())}

func getClient(projectID string, topic string) *bus.GCPClient {
	c := client
	c.ProjectID = projectID
	c.Opts = MockServerOptions
	err := c.Connect()
	if err != nil {
		return &bus.GCPClient{}
	}
	c.Instance.CreateTopic(c.Ctx, topic)
	return &c
}

func getMockServer() *grpc.ClientConn {
	srv := pstest.NewServer()
	conn, _ := grpc.Dial(srv.Addr, grpc.WithInsecure())
	return conn
}

const (
	projectID = "0cba82ff-9790-454d-b7b9-22570e7ba28c"
	prefix    = "trigger_"
)

func getTopic(t string) string {
	return fmt.Sprintf("%s%s-%s", prefix, projectID, t)
}

func TestExecTrain(t *testing.T) {
	Type := "train"
	runIDs := []string{"d948b49a-57d3-4254-994b-528509ca5963"}
	tests := []struct {
		name    string
		in      []byte
		want    *handlers.OutputPayload
		isError bool
	}{
		{
			name: "Positive",
			in:   []byte(fmt.Sprintf(`{"project_id": "%s", "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37", "pipeline_config": [{"data": {"location": {"source": "gcs://test/test.csv"}, "prep_config": {}}, "model": {"hyperparameters": {}, "version": "v1"}}]}`, projectID)),
			want: &handlers.OutputPayload{
				Errors:      []string{},
				SubmittedID: runIDs,
			},
			isError: false,
		},
		{
			name:    "Negative: json parsing error",
			in:      []byte(fmt.Sprintf(`{"project_id": "%s"`, projectID)),
			want:    &handlers.OutputPayload{},
			isError: true,
		},
		{
			name:    "Negative: validation error",
			in:      []byte(fmt.Sprintf(`{"project_id": "%s", "code_hash": "foobar", "pipeline_config": [{"data": {}, "model": {}}]}`, projectID)),
			want:    &handlers.OutputPayload{},
			isError: true,
		},
	}

	proc := handlers.Processor{
		Type:            Type,
		TopicPrefix:     prefix,
		InputJSONSchema: config.InputJSONSchemaTrain,
		Bus:             getClient(projectID, getTopic(Type)),
	}

	for _, test := range tests {
		got, err := proc.Exec(test.in)
		if !test.isError {
			if err != nil {
				t.Fatalf("[%s]: Error: %s", test.name, err)
			}
			got.SubmittedID = runIDs
			if !reflect.DeepEqual(got, test.want) {
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

func TestExecPred(t *testing.T) {
	Type := "predict"
	runIDs := []string{"37f48a01-53c7-41bc-974f-44017df781dd"}
	tests := []struct {
		name    string
		in      []byte
		want    *handlers.OutputPayload
		isError bool
	}{
		{
			name: "Positive",
			in:   []byte(fmt.Sprintf(`{"project_id": "%s", "train_id": "d948b49a-57d3-4254-994b-528509ca5963", "pipeline_config": [{"data": {"location": {"source": "gcs://test/test.csv", "destination": "gcs://test/test1.csv"}}}]}`, projectID)),
			want: &handlers.OutputPayload{
				Errors:      []string{},
				SubmittedID: runIDs,
			},
			isError: false,
		},
		{
			name: "Positive: submission error, no destination topic",
			in:   []byte(`{"project_id": "48719f69-6e3d-4ec9-a876-b8777cda74f9", "train_id": "d948b49a-57d3-4254-994b-528509ca5963", "pipeline_config": [{"data": {"location": {"source": "gcs://test/test.csv", "destination": "gcs://test/test1.csv"}}}]}`),
			want: &handlers.OutputPayload{
				Errors:      []string{`rpc error: code = NotFound desc = topic "projects/0cba82ff-9790-454d-b7b9-22570e7ba28c/topics/trigger_48719f69-6e3d-4ec9-a876-b8777cda74f9-predict". Details:nmap[data:map[location:map[destination:gcs://test/test1.csv source:gcs://test/test.csv]]]`},
				SubmittedID: []string{},
			},
			isError: false,
		},
		{
			name:    "Negative: json parsing error",
			in:      []byte(fmt.Sprintf(`{"project_id": "%s"`, projectID)),
			want:    &handlers.OutputPayload{},
			isError: true,
		},
		{
			name:    "Negative: validation error",
			in:      []byte(fmt.Sprintf(`{"project_id": "%s", "pipeline_config": [{"data": {}}]}`, projectID)),
			want:    &handlers.OutputPayload{},
			isError: true,
		},
	}

	proc := handlers.Processor{
		Type:            Type,
		TopicPrefix:     prefix,
		InputJSONSchema: config.InputJSONSchemaPredict,
		Bus:             getClient(projectID, getTopic(Type)),
	}

	for _, test := range tests {
		got, err := proc.Exec(test.in)
		if !test.isError {
			if err != nil {
				t.Fatalf("[%s]: Error: %s", test.name, err)
			}
			if test.name == "Positive" {
				got.SubmittedID = runIDs
				if !reflect.DeepEqual(got, test.want) {
					t.Fatalf("[%s]: Results don't match\nwant: %s\ngot: %s",
						test.name, mustMarshal(test.want), mustMarshal(got))
				}
			}
		} else {
			if err == nil {
				t.Fatalf("[%s]: Error wasn't thrown.", test.name)
			}
		}
	}
}
