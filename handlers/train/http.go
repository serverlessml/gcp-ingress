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
	"net/http"

	"github.com/serverlessml/gcp-ingress/handlers"
)

// HandlerPOST http handler to invoke train pipeline.
func (p *Processor) HandlerPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		handlers.HandlerError(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	inputPayload := handlers.GetRequestPayload(r.Body)

	output, err := p.Exec(inputPayload)
	if err != nil {
		handlers.HandlerError(w, err.Error(), http.StatusBadRequest)
		return
	}

	outputPayload := handlers.MustMarshal(output)
	w.WriteHeader(http.StatusAccepted)
	w.Write(outputPayload)
	return
}
