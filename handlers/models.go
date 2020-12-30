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

// Input defines the input payload to invoke pipeline.
type Input struct {
	// Model project ID (there may be several model projects within one cloud project).
	ProjectID string `json:"project_id"`
	// CodeHash is the model codebase ID.
	CodeHash string `json:"code_hash,omitempty"`
	// TrainID is the experiment's ID.
	TrainID string `json:"train_id,omitempty"`
	// Config is the ML pipeline config
	Config []interface{} `json:"pipeline_config"`
}

// OutputPayload defines output payload.
type OutputPayload struct {
	Errors      []string `json:"errors"`
	SubmittedID []string `json:"submitted_id"`
}

// OutputDistribution defines the output distribution config.
type OutputDistribution struct {
	// Topic is the message broker topic to push payload to.
	Topic string
}
