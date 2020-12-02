// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/serverlessml/gcp-ingress/processor"
)

// GetEnv extracts envvar with default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetPayload get request body's payload.
func GetPayload(requestBody io.ReadCloser) ([]byte, error) {
	payload, err := ioutil.ReadAll(requestBody)
	defer requestBody.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading payload: %s", err)
	}
	return payload, nil
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func handlerPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	payload, err := GetPayload(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := proc.Exec(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, _ := proc.Marshal(output)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

var proc processor.Processor

func main() {
	proc.TopicPrefix = GetEnv("TOPIC_PREFIX", "trigger-")

	Port := GetEnv("PORT", "8080")

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/", handlerPOST)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", Port), nil))
}
