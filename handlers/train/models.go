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

package train

// Location defines object location.
type Location struct {
	Source string `json:"source" validate:"required,startswith=fs|gcs|s3"`
}

// DataConfig defines data preparation config.
type DataConfig struct {
	Location   Location               `json:"location" validate:"required,structonly"`
	PrepConfig map[string]interface{} `json:"prep_config" validate:"required,structonly"`
}

// ModelConfig defines model definition+train config.
type ModelConfig struct {
	Hyperparameters map[string]interface{} `json:"hyperparameters" validate:"required,structonly"`
	Version         string                 `json:"version" validate:"required"`
}

// PipelineConfig defines ML pipeline config.
type PipelineConfig struct {
	// Data represents the configuration of the data preparation for an ML experiment
	Data DataConfig `json:"data" validate:"required,dive,structonly"`
	// Model represents the model setting configuration
	Model ModelConfig `json:"model" validate:"required,dive,structonly"`
}

// Input defines the input payload to invoke train pipeline.
type Input struct {
	// Model project ID (there may be several model projects within one cloud project).
	ProjectID string `json:"project_id" validate:"required,uuid4|uuid4_rfc4122"`
	// CodeHash is the model codebase ID.
	CodeHash string `json:"code_hash" validate:"required,sha1"`
	// Config is the ML pipeline config
	// it contains data preparation as well as the ML settings config
	Config []PipelineConfig `json:"pipeline_config" validate:"required,dive,structonly"`
}

// PushPayload defines the payload returned for further transition down the ML pipeline.
type PushPayload struct {
	// CodeHash is the model codebase ID.
	CodeHash string `json:"code_hash"`
	// RunID is the experiment's ID.
	RunID string `json:"run_id"`
	// Config is the ML pipeline config
	Config PipelineConfig `json:"pipeline_config"`
}
