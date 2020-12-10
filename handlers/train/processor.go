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

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/serverlessml/gcp-ingress/bus"
	"github.com/serverlessml/gcp-ingress/handlers"
)

// Processor defines processor for predict pipeline.
type Processor struct {
	// ProjectID represents cloud Project ID.
	ProjectID string
	// TopicPrefix represents prefix of the topic to post the payload to.
	TopicPrefix string
	// Message bus to disctribute messages
	Bus *bus.Client
}

// Exec run processor sequence.
func (p *Processor) Exec(data []byte) (*handlers.OutputPayload, error) {
	errs := handlers.Validate(InputJSONSchema, data)
	if errs != nil {
		return &handlers.OutputPayload{}, errs
	}

	var input Input
	err := json.Unmarshal(data, &input)
	if err != nil {
		return &handlers.OutputPayload{}, handlers.NewUnmarshallerError(err)
	}

	topic := fmt.Sprintf("%s%s-train", p.TopicPrefix, input.ProjectID)

	errorsCh := make(chan error, len(input.Config))
	runIDs := []string{}
	for _, config := range input.Config {
		runID := fmt.Sprintf("%s", uuid.NewV4())
		payloadPush, _ := json.Marshal(PushPayload{
			CodeHash: input.CodeHash,
			RunID:    runID,
			Config:   config,
		})
		go p.Bus.PushRoutine(payloadPush, topic, errorsCh)
		runIDs = append(runIDs, runID)
	}

	pushErrors := []string{}
	outputRunIDs := []string{}
	for i, config := range input.Config {
		err := <-errorsCh
		if err != nil {
			e := handlers.ErrorPush{
				Message: err.Error(),
				Details: config,
			}
			pushErrors = append(pushErrors, e.Error())
		} else {
			outputRunIDs = append(outputRunIDs, runIDs[i])
		}
	}

	return &handlers.OutputPayload{
		Errors:      pushErrors,
		SubmittedID: outputRunIDs,
	}, nil
}
