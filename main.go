// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

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
func GetRequestPayload(requestBody io.ReadCloser) ([]byte, error) {
	payload, err := ioutil.ReadAll(requestBody)
	defer requestBody.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading payload: %s", err)
	}
	return payload, nil
}

// MustMarshal performs json.Marshal
func MustMarshal(obj interface{}) []byte {
	out, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	return out
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
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

	inputPayload, err := GetRequestPayload(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

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
)

func main() {
	proc.TopicPrefix = GetEnv("TOPIC_PREFIX", "trigger_")
	pubsubClient.ProjectID = GetEnv("PROJECT_ID", "project")

	err := pubsubClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	Port := GetEnv("PORT", "8080")

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/", handlerPOST)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port), nil))
}
