#! /bin/bash

curl -isX GET "http://0.0.0.0:8080/status" | grep "200 OK" | if [ $(wc -l) -eq 0 ]; then echo "FAIL: Health check unreachable"; exit 1; fi

curl -isX POST \
    -d '{"project_id": "0cba82ff-9790-454d-b7b9-22570e7ba28c", "code_hash": "8c2f3d3c5dd853231c7429b099347d13c8bb2c37","pipeline_config": [{"data": {"foo": 1},"model": {"foo": 1}},{"data": {"foo": 2},"model": {"bar": 2}}]}' \
    "http://0.0.0.0:8080/" | grep "202 Accepted" | if [ $(wc -l) -eq 0 ]; then echo "FAIL:Faulty main endpoint"; exit 1; fi

echo "OK"
