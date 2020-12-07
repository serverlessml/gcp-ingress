// Dmitry Kisler © 2020-present
// www.dkisler.com <admin@dkisler.com>

package bus

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

// Client defines the client to communicate with the message bus.
type Client struct {
	// GCP Project ID
	ProjectID string
	Ctx       context.Context
	Instance  *pubsub.Client
}

// Connect establishes connector to the message broker.
func (c *Client) Connect(opts ...option.ClientOption) error {
	var err error
	c.Ctx = context.Background()
	c.Instance, err = pubsub.NewClient(c.Ctx, c.ProjectID, opts...)
	return err
}

// Push pushes the message to a topic.
func (c *Client) Push(payload []byte, topic string) error {
	t := c.Instance.Topic(topic)
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(c.Ctx, &pubsub.Message{Data: payload})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err := result.Get(c.Ctx)
	return err
}

// PushRoutine marshal the output payload object.
func (c *Client) PushRoutine(payload []byte, topic string, ch chan error) {
	err := c.Push(payload, topic)
	ch <- err
}
