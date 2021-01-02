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

package handlers

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	bus "github.com/serverlessml/platform/gcp/bus"
)

// Processor defines processor for predict pipeline.
type Processor struct {
	// Type defines the execution type, either train, or predict
	Type string
	// TopicPrefix represents prefix of the topic to post the payload to.
	TopicPrefix string
	// InputJSONSchema defines the jsonschema of the input payload.
	InputJSONSchema string
	// Message bus to disctribute messages
	Bus *bus.Client
}

// Exec run processor sequence.
func (p *Processor) Exec(data []byte) (*OutputPayload, error) {
	errs := Validate(p.InputJSONSchema, data)
	if errs != nil {
		return nil, errs
	}

	var input Input
	json.Unmarshal(data, &input)

	errorsCh, runIDs := p.distubuteData(input)

	return p.formatOutput(input.Config, errorsCh, runIDs), nil
}

// distubuteData distributes data further down pipeline.
// any output interfaces can be plugged in here, e.g. pubsub, kafka, db
func (p *Processor) distubuteData(input Input) (chan error, []string) {
	topic := fmt.Sprintf("%s%s-%s", p.TopicPrefix, input.ProjectID, p.Type)
	errorsCh := make(chan error, len(input.Config))
	runIDs := []string{}
	for _, config := range input.Config {
		runID := fmt.Sprintf("%s", uuid.NewV4())

		payload := map[string]interface{}{
			"run_id": runID,
			"config": config,
		}
		switch p.Type {
		case "train":
			payload["code_hash"] = input.CodeHash
		case "predict":
			payload["train_id"] = input.TrainID
		}

		go p.Bus.PushRoutine(MustMarshal(payload), topic, errorsCh)

		runIDs = append(runIDs, runID)
	}
	return errorsCh, runIDs
}

// formatOutput formats output of the main processor's method.
func (p *Processor) formatOutput(configs []interface{}, errorsCh chan error, runIDs []string) *OutputPayload {
	pushErrors := []string{}
	outputRunIDs := []string{}
	for i, config := range configs {
		err := <-errorsCh
		if err != nil {
			e := ErrorPush{
				Message: err.Error(),
				Details: config,
			}
			pushErrors = append(pushErrors, e.Error())
		} else {
			outputRunIDs = append(outputRunIDs, runIDs[i])
		}
	}
	return &OutputPayload{
		Errors:      pushErrors,
		SubmittedID: outputRunIDs,
	}
}
