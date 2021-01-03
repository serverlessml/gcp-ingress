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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/serverlessml/ingress/config"
	"github.com/serverlessml/ingress/handlers"
)

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
		got := handlers.GetRequestPayload(test.in)
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
		got := handlers.MustMarshal(test.in)

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

func TestHandlerStatus(t *testing.T) {
	tests := []struct {
		name string
		in   *http.Request
		want int
	}{
		{
			name: "Positive",
			in:   httptest.NewRequest("GET", "http://0.0.0.0:8080", nil),
			want: http.StatusOK,
		},
		{
			name: "Negative: wrong request type",
			in:   httptest.NewRequest("POST", "http://0.0.0.0:8080", nil),
			want: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		handlers.HandlerStatus(w, test.in)
		got := w.Result().StatusCode
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %d\ngot: %d",
				test.name, test.want, got)
		}
	}
}

func TestHandlerError(t *testing.T) {
	type args struct {
		errMsg []error
		status int
	}
	tests := []struct {
		in   *args
		want string
	}{
		{
			in: &args{
				errMsg: []error{fmt.Errorf("foobar")},
				status: 404,
			},
			want: `{"errors":["foobar"],"submitted_id":[]}`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		handlers.HandlerError(w, test.in.errMsg, test.in.status)
		got := w.Result()
		gotStatusCode := got.StatusCode

		if gotStatusCode != test.in.status {
			t.Fatalf("Results don't match\nwant: %d\ngot: %d", test.in.status, gotStatusCode)
		}

		body, _ := ioutil.ReadAll(got.Body)
		gotResp := string(body)
		if gotResp != test.want {
			t.Fatalf("Results don't match\nwant: %s\ngot: %s", test.want, gotResp)
		}
	}
}

func TestHandlerPOSTTrain(t *testing.T) {
	projectID := "0cba82ff-9790-454d-b7b9-22570e7ba28c"
	prefix := "trigger_"
	Type := "train"
	proc := handlers.Processor{
		Type:            Type,
		TopicPrefix:     prefix,
		InputJSONSchema: config.InputJSONSchemaTrain,
		Bus:             getClient(projectID, getTopic(Type)),
	}

	tests := []struct {
		name string
		in   *http.Request
		want int
	}{
		{
			name: "Positive",
			in: httptest.NewRequest("POST", "/train", strings.NewReader(fmt.Sprintf(`{
  "project_id": "%s",
  "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
  "pipeline_config": [
	{
	  "data": {
		"location": {
		  "source": "gcs://test/train.csv"
		},
		"prep_config": {}
	  },
	  "model": {
		"hyperparameters": {},
		"version": "v1"
	  }
	}
  ]
}`, projectID))),
			want: http.StatusAccepted,
		},
		{
			name: "Negative: parsing error",
			in: httptest.NewRequest("POST", "/train", strings.NewReader(fmt.Sprintf(`{
  "project_id": "%s",
  "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
  "pipeline_config": [
	{
	  "data": {
		"location": {
		  "source": "gcs://test/train.csv"
		},
		"prep_config": {}
	  },
	  "model": {
		"hyperparameters": {},
		"version": "v1"
	  }
	}
  ]`, projectID))),
			want: http.StatusBadRequest,
		},
		{
			name: "Negative: validation error",
			in: httptest.NewRequest("POST", "/train", strings.NewReader(fmt.Sprintf(`{
  "project_id": "%s",
  "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
  "pipeline_config": [
	{
	  "data": {
		"location": {
		  "source": "gcs://test/train.csv"
		},
		"prep_config": {}
	  },
	  "model": {
		"hyperparameters": {}
	  }
	}
  ]
}`, projectID))),
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "Negative: wrong request type",
			in:   httptest.NewRequest("GET", "/train", nil),
			want: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		handlers.HandlerPOST(w, test.in, &proc)
		got := w.Result().StatusCode
		if got != test.want {
			t.Fatalf("[%s]: Results don't match\nwant: %d\ngot: %d",
				test.name, test.want, got)
		}
	}
}
