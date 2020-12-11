// // Copyright 2020 dkisler.com Dmitry Kisler
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// // EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// // OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND
// // NONINFRINGEMENT. IN NO EVENT WILL THE LICENSOR OR OTHER CONTRIBUTORS
// // BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN
// // ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF, OR IN
// // CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// //
// // See the License for the specific language governing permissions and
// // limitations under the License.

package config

// InputJSONSchemaTrain defines json schema for input payload to invoke train pipeline.
const InputJSONSchemaTrain string = `{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "definitions": {
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
                "source"
            ],
            "properties": {
                "source": {
                    "$ref": "#/definitions/path"
                }
            }
        },
        "data_config": {
            "type": "object",
            "description": "Data prep config.",
            "additionalProperties": false,
            "required": [
                "location",
                "prep_config"
            ],
            "properties": {
                "location": {
                    "$ref": "#/definitions/data_location"
                },
                "prep_config": {
                    "type": "object",
                    "description": "Config to prepare data set for model training."
                }
            }
        },
        "model_config": {
            "type": "object",
            "description": "Model train config.",
            "additionalProperties": false,
            "required": [
                "hyperparameters",
                "version"
            ],
            "properties": {
                "hyperparameters": {
                    "type": "object",
                    "description": "Model's hyperparameters configuration."
                },
                "version": {
                    "type": "string",
                    "description": "Model's version name.",
                    "pattern": "^[a-zA-Z0-9_\\-.|]{1,40}$"
                }
            }
        },
        "item": {
            "type": "object",
            "description": "Data prep config.",
            "additionalProperties": false,
            "required": [
                "data",
                "model"
            ],
            "properties": {
                "data": {
                    "$ref": "#/definitions/data_config"
                },
                "model": {
                    "$ref": "#/definitions/model_config"
                }
            }
        }
    },
    "type": "object",
    "title": "Train trigger schema",
    "additionalProperties": false,
    "required": [
        "project_id",
        "code_hash",
        "pipeline_config"
    ],
    "properties": {
        "project_id": {
            "description": "Modelling project ID.",
            "type": "string",
            "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
            "examples": [
                "0cba82ff-9790-454d-b7b9-22570e7ba28c"
            ]
        },
        "code_hash": {
            "description": "Codebase (git commit) SHA1 hash value.",
            "type": "string",
            "pattern": "^[a-fA-F0-9]{40}$",
            "examples": [
                "8c2f3d3c5dd853231c7429b099347d13c8bb2c37"
            ]
        },
        "pipeline_config": {
            "type": "array",
            "description": "ML pipeline configuration for training.",
            "items": {
                "$ref": "#/definitions/item"
            }
        }
    }
}`
