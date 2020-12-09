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

package predict

// Location defines object location.
type Location struct {
	Source      string `json:"source" validate:"required,startswith=fs|gcs|s3"`
	Destination string `json:"destination" validate:"required,startswith=fs|gcs|s3"`
}

// DataConfig defines data preparation config.
type DataConfig struct {
	Location Location `json:"location" validate:"required,structonly"`
}

// PipelineConfig defines ML pipeline config.
type PipelineConfig struct {
	// Data represents the configuration of the data preparation for an ML experiment
	Data DataConfig `json:"data" validate:"required,dive,structonly"`
}

// Input defines the input payload to invoke pipeline for prediction.
type Input struct {
	// Model project ID (there may be several model projects within one cloud project).
	ProjectID string `json:"project_id" validate:"required,uuid4|uuid4_rfc4122"`
	// RunID is the experiment's ID.
	TrainID string `json:"train_id" validate:"required,uuid4|uuid4_rfc4122"`
	// Config is the ML pipeline config
	Config []PipelineConfig `json:"pipeline_config" validate:"required,dive,structonly"`
}

// PushPayload defines the payload returned for further transition down the ML pipeline.
type PushPayload struct {
	// RunID is the prediction run ID.
	RunID string `json:"run_id"`
	// RunID is the experiment's ID.
	TrainID string `json:"train_id" validate:"required,uuid4|uuid4_rfc4122"`
	// Config is the ML pipeline config
	Config PipelineConfig `json:"pipeline_config"`
}
