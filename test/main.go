package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	a := `{
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
}`

	out, _ := json.Marshal(a)
	fmt.Print(out)
}
