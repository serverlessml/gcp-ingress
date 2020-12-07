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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/serverlessml/gcp-ingress/bus"
	"github.com/serverlessml/gcp-ingress/processor"
)

// GetEnv extracts envvar with default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetRequestPayload get request body's payload.
func GetRequestPayload(requestBody io.ReadCloser) []byte {
	payload, _ := ioutil.ReadAll(requestBody)
	defer requestBody.Close()
	return payload
}

// MustMarshal performs json.Marshal
func MustMarshal(obj interface{}) []byte {
	out, _ := json.Marshal(obj)
	return out
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

// errorOutput defines output error.
type errorOutput struct {
	// contains error message
	Message string `json:"message"`
	// contains pipeline config
	PipelineConfig processor.PipelineConfig `json:"pipeline_config"`
}

// OutputPayload defines output payload.
type OutputPayload struct {
	Errors      []errorOutput `json:"errors"`
	SubmittedID []string      `json:"submitted_id"`
}

// runner defines the main routine.
func runner(data []byte) (OutputPayload, error) {
	outputProc, err := proc.Exec(data)
	if err != nil {
		return OutputPayload{}, err
	}
	outputProcPayload := &outputProc.Payload

	errorsCh := make(chan error, len(*outputProcPayload))
	for _, payload := range *outputProcPayload {
		payloadProcOutput, _ := json.Marshal(payload)
		go pubsubClient.PushRoutine(payloadProcOutput, outputProc.Distribution.Topic, errorsCh)
	}

	outputErrors := []errorOutput{}
	outputRunIDs := []string{}
	for _, item := range *outputProcPayload {
		err := <-errorsCh
		if err != nil {
			outputErrors = append(outputErrors, errorOutput{
				Message:        err.Error(),
				PipelineConfig: item.Config,
			})
		} else {
			outputRunIDs = append(outputRunIDs, item.RunID)
		}
	}

	return OutputPayload{
		Errors:      outputErrors,
		SubmittedID: outputRunIDs,
	}, nil
}

func errorResponse(w http.ResponseWriter, errMsg string, status int) {
	w.WriteHeader(status)
	w.Write(MustMarshal(OutputPayload{
		Errors: []errorOutput{{
			Message:        errMsg,
			PipelineConfig: processor.PipelineConfig{},
		}},
		SubmittedID: []string{},
	}))
	return
}

func handlerPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorResponse(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	inputPayload := GetRequestPayload(r.Body)

	output, err := runner(inputPayload)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	outputPayload := MustMarshal(output)
	w.WriteHeader(http.StatusAccepted)
	w.Write(outputPayload)
	return
}

var (
	proc         processor.Processor
	pubsubClient bus.Client
	httpSrv      http.Server
)

func main() {
	proc.TopicPrefix = GetEnv("TOPIC_PREFIX", "trigger_")
	pubsubClient.ProjectID = GetEnv("PROJECT_ID", "project")
	httpSrv.Addr = fmt.Sprintf(":%s", GetEnv("PORT", "8080"))

	err := pubsubClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/status", handlerStatus)
	http.HandleFunc("/", handlerPOST)
	go func() {
		if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
}
