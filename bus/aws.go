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

package bus

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AWSClient defines the client to communicate with the SNS.
type AWSClient struct {
	// AWS project ID
	ProjectID string
	// AWS region
	Region   string
	Session  *session.Session
	Instance *sns.SNS
}

// getSession establishes the connection session.
func (c *AWSClient) getSession() error {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(c.Region),
	})
	if err != nil {
		return nil
	}
	c.Session = ses
	return nil
}

// getIdentity fetches the project ID using AWS STS.
func (c *AWSClient) getProjectID() error {
	stsClient := sts.New(c.Session)
	identity, err := stsClient.GetCallerIdentity(nil)
	if err != nil {
		return err
	}
	c.ProjectID = *identity.Account
	return nil
}

// Connect establishes connector to the message broker.
func (c *AWSClient) Connect() error {
	err := c.getSession()
	err = c.getProjectID()
	c.Instance = sns.New(c.Session)
	return err
}

// getTopicArn fetches the topic ARN based on it's name.
func (c *AWSClient) getTopicArn(topic string) string {
	return fmt.Sprintf("arn:aws:sns:%s:%s:%s", c.Region, c.ProjectID, topic)
}

// Push pushes the message to a topic.
func (c *AWSClient) Push(payload []byte, topic string) error {
	t := c.getTopicArn(topic)

	_, err := c.Instance.Publish(&sns.PublishInput{
		Message:  aws.String(string(payload)),
		TopicArn: &t,
	})

	return err
}
