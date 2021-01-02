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

// BusClient defines the client to push message to message bus.
type BusClient interface {
	Connect() error
	Push([]byte, string) error
}

// NewBusClient creates a new client instance.
// func NewBusClient(platform string, conf map[string]string) (BusClient, error) {
// 	var client BusClient
// 	var err error
// 	switch platform {
// 	case "aws":
// 		if conf["region"] == "" {
// 			return nil, Error{
// 				Type:    "init bus client",
// 				Message: fmt.Sprintf("%s bus client requires 'region' config.", platform),
// 			}
// 		}
// 		client = &bus.AWSClient{Region: conf["region"]}
// 		err = client.Connect()
// 	case "gcp":
// 		if conf["project_id"] == "" {
// 			return nil, Error{
// 				Type:    "init bus client",
// 				Message: fmt.Sprintf("%s bus client requires 'project_id' config.", platform),
// 			}
// 		}
// 		client = &bus.GCPClient{ProjectID: conf["project_id"]}
// 		err = client.Connect()
// 	default:
// 		return nil, Error{
// 			Type:    "init bus client",
// 			Message: fmt.Sprintf("%s platform not implemented.", platform),
// 		}
// 	}
// 	return client, err
// }

// PushRoutine pushes the message to a topic for async go-routines.
func PushRoutine(c BusClient, payload []byte, topic string, ch chan error) {
	err := c.Push(payload, topic)
	ch <- err
}
