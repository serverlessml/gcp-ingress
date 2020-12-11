# Pipeline Ingress: GCP

[![Go Report Card](https://goreportcard.com/badge/github.com/serverlessml/gcp-ingress)](https://goreportcard.com/report/github.com/serverlessml/gcp-ingress) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/741b0eb31b494469a2cedb2046fe60fb)](https://www.codacy.com/gh/serverlessml/gcp-ingress/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=serverlessml/gcp-ingress&amp;utm_campaign=Badge_Grade) [![Code Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://img.shields.io/badge/coverage-100%25-brightgreen)

The ingress service to invoke ML pipeline. A web-server with two end-points:

`GET: /status` -> status check

`POST: /`      -> pipeline invocation trigger request with the metadata payload

## Modus Operandi

1. Validate input payload:

```js
{
    "id": UUID4,
    "config": {
        "model": Object,
        "data": Object,
    }
}
```

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

### Tests

#### Health check

```bash
curl -iX GET "http://0.0.0.0:8080/status"
```

Expected output:

```bash
HTTP/1.1 200 OK
Date: YOUR CURRENT DATE/TIME
Content-Length: 0
```

#### Integration test

1. Create the PubSub topic called `trigger_0cba82ff-9790-454d-b7b9-22570e7ba28c`
2. Execute `make PROJECT_ID=<YOUR_GCO_PROJECT_ID> run`
3. Run

```bash
curl -iX POST "http://0.0.0.0:8080/" -d '{"id": "0cba82ff-9790-454d-b7b9-22570e7ba28c", "config": {"data": {}, "model": {}}}'
```

Expected output:
```bash
HTTP/1.1 202 Accepted
Date: YOUR CURRENT DATE/TIME
Content-Length: 0
```
