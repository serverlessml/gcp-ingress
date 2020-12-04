// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package processor

import (
	"encoding/json"
	"fmt"

	"github.com/serverlessml/gcp-ingress/validator"
)

// config defines ML pipeline config.
type config struct {
	// Data represents the configuration of the data preparation for an ML experiment
	Data map[string]interface{} `json:"data"`
	// Model represents the model setting configuration
	Model map[string]interface{} `json:"model"`
}

// input defines the input payload.
type input struct {
	// ProjectID is the project ID
	ProjectID string `json:"project_id" validate:"required,uuid4|uuid_rfc4122"`
	// CodeHash is the model codebase ID
	CodeHash string `json:"code_hash" validate:"required,sha1"`
	// Config is the ML pipeline config
	// it contains data preparation as well as the ML settings config
	Config []config `json:"config" validate:"required"`
}

// outputDistribution defines the output distribution config.
type outputDistribution struct {
	// Topic is the message broker topic to push payload to.
	Topic string
}

// Output defines the output object.
type Output struct {
	// Payload represents the output config payload.
	Payload []config `json:"config"`
	// Distribution defines the payload distribution config.
	Distribution outputDistribution
}

// Processor defines processor.
type Processor struct {
	// TopicPrefix represents prefix of the topic to post the payload to.
	TopicPrefix string
}

// readInput reads the input data content.
func readInput(data []byte) (*input, error) {
	var inpt input
	err := json.Unmarshal(data, &inpt)
	return &inpt, err
}

// validateInput validates the input.
func validateInput(input *input) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err == nil {
		return nil
	}
	validationErrors := validator.GetValidationErrors(err)
	if len(validationErrors) == 0 {
		return nil
	}
	return fmt.Errorf(fmt.Sprintln(validationErrors))
}

// Exec run processor sequence.
func (p *Processor) Exec(data []byte) (*Output, error) {
	input, err := readInput(data)
	if err != nil {
		return &Output{}, fmt.Errorf("Input reading error: %s", err)
	}

	err = validateInput(input)
	if err != nil {
		return &Output{}, fmt.Errorf("Input validation error: %s", err)
	}

	return &Output{
		Payload: input.Config,
		Distribution: outputDistribution{
			Topic: fmt.Sprintf("%s%s", p.TopicPrefix, input.ProjectID),
		},
	}, err
}
