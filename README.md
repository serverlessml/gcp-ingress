# Pipeline Ingress: GCP

[![Go Report Card](https://goreportcard.com/badge/github.com/serverlessml/gcp-ingress)](https://goreportcard.com/report/github.com/serverlessml/gcp-ingress) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/741b0eb31b494469a2cedb2046fe60fb)](https://www.codacy.com/gh/serverlessml/gcp-ingress/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=serverlessml/gcp-ingress&amp;utm_campaign=Badge_Grade) [![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://img.shields.io/badge/coverage-100%25-brightgreen)

The ingress service to invoke ML pipeline. A web-server with two end-points:

`GET: /status`      -> status check

`POST: /train`      -> <strong><em>train</em></strong> pipeline invocation trigger request with the metadata payload

`POST: /predict`    -> <strong><em>predict</em></strong> pipeline invocation trigger request with the metadata payload

## Modus Operandi

1. Validate input payload:
2. Push the payload to a GCP PubSub topic
3. Return 202 as response in case of success

## How to run
!Note! It requires an existing GCP account with activated pubsub API, service account (SA) with PubSub Write permissions and a pubsub topic. Once SA is created, generate and download the access key to `${PATH_TO_SERVICE_ACCOUNT_KEY}/key-pubsub.json`.

Execute:

```bash
make build
```

afterwards:

```bash
make PROJECT_ID=<YOUR_GCO_PROJECT_ID> run
```

or alternatively, for the sake of testing, run

```bash
make PROJECT_ID=<YOUR_GCO_PROJECT_ID> test-run
```

### HTTP Response Codes
|Endpoint|Method|HTTP Status Code|Comment|
|:-|:-:|-:|--|--|
|/status|GET|200|-|
|/status|POST,PUT,PATCH,DELETE|405|Not supported methods|
|/train<br>/predict|POST|202|Request accepted|
|/train<br>/predict|POST|400|Faulty JSON submitted with request|
|/train<br>/predict|GET,PUT,PATCH,DELETE|405|Not supported methods|
|/train<br>/predict|POST|422|Submitted JSON doesn't pass validation/comply with the input schema|

#### Input Schema

Input JSON schemas per endpoint can be found here:

- [/train](./config/schema_train.go)

- [/predict](./config/schema_predict.go)

#### Response Schema

POST requests sent to `/train` and `/predict` lead to responses of the following structure:

```json
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "definitions": {
        "uuid4": {
            "oneOf": [
                {
                    "type": "string",
                    "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
                },
                {
                    "type": "null"
                }
            ]
        }
    },
    "type": "object",
    "additionalProperties": false,
    "required": [
        "errors",
        "submitted_id"
    ],
    "properties": {
        "errors": {
            "type": "array",
            "items": {
                "oneOf": [
                    {
                        "type": "string"
                    },
                    {
                        "type": "null"
                    }
                ]
            }
        },
        "submitted_id": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/uuid4"
            }
        }
    }
}
```

### Tests

#### Health check

```bash
curl -iX GET "http://0.0.0.0:8080/status"
```

Expected output:

```bash
HTTP/1.1 200 OK
Access-Control-Allow-Methods: GET
Date: YOUR CURRENT DATE/TIME
Content-Length: 0
```

#### Integration test

1. Create the PubSub topic called `trigger_0cba82ff-9790-454d-b7b9-22570e7ba28c`
2. Execute `make PROJECT_ID=<YOUR_GCO_PROJECT_ID> run`
3. Run

- to trigger `train` pipeline:

```bash
curl -H "Accept: application/json" -H "Content-Type: application/json" \
    -iX POST "http://0.0.0.0:8080/train" -d '{
  "project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c",
  "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37",
  "pipeline_config": [
    {
      "data": {
        "location": {
          "source": "gcs://test/train.csv"
        },
        "prep_config": {}
      },
      "model": {
        "hyperparameters": {},
        "version": "v1"
      }
    }
  ]
}'
```

Expected output:
```bash
HTTP/1.1 202 Accepted
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding
Access-Control-Allow-Methods: POST
Content-Type: application/json
Date: YOUR CURRENT DATE/TIME
Content-Length: 69

{"errors":[],"submitted_id":["b441141f-bce1-4552-9458-999d2b8f6fda"]}
```

- to trigger `predict` pipeline:

```bash
curl -H "Accept: application/json" -H "Content-Type: application/json" \
    -iX POST "http://0.0.0.0:8080/predict" -d '{
  "project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c",
  "train_id": "b441141f-bce1-4552-9458-999d2b8f6fda",
  "pipeline_config": [
    {
      "data": {
        "location": {
          "source": "gcs://test/train.csv",
          "destination": "gcs://prediction/train.csv"
        }
      }
    }
  ]
}'
```

Expected output:
```bash
HTTP/1.1 202 Accepted
Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding
Access-Control-Allow-Methods: POST
Content-Type: application/json
Date: YOUR CURRENT DATE/TIME
Content-Length: 69

{"errors":[],"submitted_id":["3aab10b8-a42b-4938-aea8-fd398d2c2f01"]}
```

In the return payload, the `submitted_id` is a list UUID4 index, it must be different from the ones above.
