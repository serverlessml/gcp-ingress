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
// +build !test

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/serverlessml/gcp-ingress/bus"
	"github.com/serverlessml/gcp-ingress/config"
	"github.com/serverlessml/gcp-ingress/handlers"
)

// GetEnv extracts envvar with default value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	projectID := GetEnv("PROJECT_ID", "project")
	topicPrefix := GetEnv("TOPIC_PREFIX", "trigger_")

	pubsubClient := &bus.Client{ProjectID: projectID}
	err := pubsubClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	procTrain := &handlers.Processor{
		Type:            "train",
		TopicPrefix:     topicPrefix,
		InputJSONSchema: config.InputJSONSchemaTrain,
		Bus:             pubsubClient,
	}

	procPredict := &handlers.Processor{
		Type:            "predict",
		TopicPrefix:     topicPrefix,
		InputJSONSchema: config.InputJSONSchemaPredict,
		Bus:             pubsubClient,
	}

	http.HandleFunc("/status", handlers.HandlerStatus)
	http.HandleFunc("/train", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerPOST(w, r, procTrain)
	})
	http.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerPOST(w, r, procPredict)
	})

	log.Fatalf("ListenAndServe(): %v",
		http.ListenAndServe(fmt.Sprintf(":%s", GetEnv("PORT", "8080")), nil),
	)
}
