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
