package bus_test

import (
	"testing"

	"github.com/serverlessml/gcp-ingress/bus"
)

var client bus.Client

func TestConnect(t *testing.T) {
	c := client
	c.ProjectID = "test"
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
}
