// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package processor

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/serverlessml/gcp-ingress/validator"
)

// PipelineConfig defines ML pipeline config.
type PipelineConfig struct {
	// Data represents the configuration of the data preparation for an ML experiment
	Data map[string]interface{} `json:"data" validate:"required"`
	// Model represents the model setting configuration
	Model map[string]interface{} `json:"model" validate:"required"`
}

// Input defines the input payload.
type Input struct {
	// ProjectID is the project ID.
	ProjectID string `json:"project_id" validate:"required,uuid4|uuid_rfc4122"`
	// CodeHash is the model codebase ID.
	CodeHash string `json:"code_hash" validate:"required,sha1"`
	// Config is the ML pipeline config
	// it contains data preparation as well as the ML settings config
	Config []PipelineConfig `json:"pipeline_config" validate:"required,dive"`
}

// OutputPayload defines the payload returned for further transition down the ML pipeline.
type OutputPayload struct {
	// Config is the ML pipeline config
	Config PipelineConfig `json:"pipeline_config" validate:"required,dive"`
	// CodeHash is the model codebase ID.
	CodeHash string `json:"code_hash" validate:"required,sha1"`
	// RunID is the experiment's ID.
	RunID string `json:"run_id" validate:"required,uuid4"`
}

// OutputDistribution defines the output distribution config.
type OutputDistribution struct {
	// Topic is the message broker topic to push payload to.
	Topic string
}

// Output defines the output object.
type Output struct {
	// Payload represents the output config payload.
	Payload []OutputPayload
	// Distribution defines the payload distribution config.
	Distribution OutputDistribution
}

// Processor defines processor.
type Processor struct {
	// TopicPrefix represents prefix of the topic to post the payload to.
	TopicPrefix string
}

// readInput reads the input data content.
func readInput(data []byte) (*Input, error) {
	var inpt Input
	err := json.Unmarshal(data, &inpt)
	return &inpt, err
}

// validateInput validates the input.
func validateInput(input *Input) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err == nil {
		return nil
	}
	validationErrors := validator.GetValidationErrors(err)
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

	output := []OutputPayload{}
	for _, config := range input.Config {
		output = append(output, OutputPayload{
			Config:   config,
			CodeHash: input.CodeHash,
			RunID:    fmt.Sprintf("%s", uuid.NewV4()),
		})
	}

	return &Output{
		Payload: output,
		Distribution: OutputDistribution{
			Topic: fmt.Sprintf("%s%s", p.TopicPrefix, input.ProjectID),
		},
	}, nil
}
