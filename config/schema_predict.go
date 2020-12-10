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

package config

// InputJSONSchemaPredict defines the json schema for input payload to invoke prediction pipeline
const InputJSONSchemaPredict string = `{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "definitions": {
        "uuid4": {
            "type": "string",
            "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
        },
        "path": {
            "type": "string",
            "description": "Path to a data object.",
            "oneOf": [
                {
                    "description": "File on file system.",
                    "pattern": "^fs://.*?$"
                },
                {
                    "description": "Object in a GCS (GCP) bucket.",
                    "pattern": "^gcs://[a-zA-Z0-9_\\-.]{3,63}/.*?$"
                },
                {
                    "description": "Object in a s3 bucket.",
                    "pattern": "^s3://[a-zA-Z0-9-.]{3,63}/.*?$"
                }
            ]
        },
        "data_location": {
            "type": "object",
            "description": "Data location.",
            "additionalProperties": false,
            "required": [
                "source",
                "destination"
            ],
            "properties": {
                "source": {
                    "$ref": "#/definitions/path"
                },
                "destination": {
                    "$ref": "#/definitions/path"
                }
            }
        },
        "data_config": {
            "type": "object",
            "description": "Data prep config.",
            "additionalProperties": false,
            "required": [
                "location"
            ],
            "properties": {
                "location": {
                    "$ref": "#/definitions/data_location"
                }
            }
        },
        "pipeline_config_item": {
            "type": "object",
            "description": "Data prep config.",
            "additionalProperties": false,
            "required": [
                "data"
            ],
            "properties": {
                "data": {
                    "$ref": "#/definitions/data_config"
                }
            }
        }
    },
    "type": "object",
    "title": "Prediction trigger schema",
    "additionalProperties": false,
    "required": [
        "project_id",
        "train_id",
        "pipeline_config"
    ],
    "properties": {
        "project_id": {
            "description": "Modelling project ID.",
            "$ref": "#/definitions/uuid4"
        },
        "train_id": {
            "description": "Train experiment/run ID.",
            "$ref": "#/definitions/uuid4"
        },
        "pipeline_config": {
            "type": "array",
            "description": "ML pipeline configuration for prediction.",
            "items": {
                "$ref": "#/definitions/pipeline_config_item"
            }
        }
    }
}`
