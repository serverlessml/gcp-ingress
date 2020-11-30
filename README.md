# Pipeline Ingress: GCP

The ingress service to invoke ML pipeline. A web-server with two end-points:

`GET: /status` -> status check

`POST: /`      -> pipeline invocation trigger request with the metadata payload

## Modus Operandi

1. Validate input payload:

```js
{
    "id": UUID4,
    "model_config": Object,
    "data_config": Object,
}
```

2. Push the payload to a GCP PubSub topic
3. Return 200 as response in case of success
