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

// PushRoutine pushes the message to a topic for async go-routines.
func PushRoutine(c BusClient, payload []byte, topic string, ch chan error) {
	err := c.Push(payload, topic)
	ch <- err
}
