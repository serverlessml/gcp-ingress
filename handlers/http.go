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
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

// MustMarshal performs json.Marshal
func MustMarshal(obj interface{}) []byte {
	out, _ := json.Marshal(obj)
	return out
}

// GetRequestPayload get request body's payload.
func GetRequestPayload(requestBody io.ReadCloser) []byte {
	payload, _ := ioutil.ReadAll(requestBody)
	defer requestBody.Close()
	return payload
}

// HandlerError http error handler.
func HandlerError(w http.ResponseWriter, errs []error, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(MustMarshal(OutputPayload{
		Errors:      ErrorArray(errs),
		SubmittedID: []string{},
	}))
	return
}

// HandlerStatus http status handler.
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.WriteHeader(http.StatusOK)
	return
}

// HandlerPOST http handler to invoke train pipeline.
func HandlerPOST(w http.ResponseWriter, r *http.Request, p *Processor) {
	if r.Method != "POST" {
		HandlerError(w, []error{errors.New("Method is not supported")}, http.StatusMethodNotAllowed)
		return
	}

	inputPayload := GetRequestPayload(r.Body)

	output, err := p.Exec(inputPayload)
	if err != nil {
		var status int
		switch eType := (err.(Error)).Type; eType {
		case "parsing":
			status = http.StatusBadRequest
		case "validation":
			status = http.StatusUnprocessableEntity
		}
		HandlerError(w, []error{err}, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
	w.WriteHeader(http.StatusAccepted)
	w.Write(MustMarshal(output))
	return
}
