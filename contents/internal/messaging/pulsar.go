package messaging

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/pulsar-client-go/pulsar"
)

var (
	pulsarClient pulsar.Client
	producer     pulsar.Producer
	consumer     pulsar.Consumer
)

// Init initialises the Pulsar client.  Connection parameters are read from
// PAO-injected environment variables (camelCase secret keys → UPPER_SNAKE):
//
//	MESSAGING_BROKER_URL               — pulsar broker URL (e.g. pulsar://broker:6650)
//	MESSAGING_TOPIC            — topic name
//	MESSAGING_JWT_TOKEN                 — JWT bearer token (empty → no auth)
//	MESSAGING_SUBSCRIPTION_NAME — consumer subscription name (consume access only)
func Init(brokerURL, topic, jwtToken, subscriptionName string) error {
	opts := pulsar.ClientOptions{URL: brokerURL}
	if jwtToken != "" {
		opts.Authentication = pulsar.NewAuthenticationToken(jwtToken)
	}

	c, err := pulsar.NewClient(opts)
	if err != nil {
		return fmt.Errorf("messaging: client: %w", err)
	}

	access := os.Getenv("MESSAGING_ACCESS")
	if access == "consume" {
		sub := subscriptionName
		if sub == "" {
			sub = topic + "-sub"
		}
		cons, err := c.Subscribe(pulsar.ConsumerOptions{
			Topic:            topic,
			SubscriptionName: sub,
			Type:             pulsar.Shared,
		})
		if err != nil {
			c.Close()
			return fmt.Errorf("messaging: consumer: %w", err)
		}
		consumer = cons
	} else {
		p, err := c.CreateProducer(pulsar.ProducerOptions{Topic: topic})
		if err != nil {
			c.Close()
			return fmt.Errorf("messaging: producer: %w", err)
		}
		producer = p
	}

	pulsarClient = c
	return nil
}

func Close() {
	if producer != nil {
		producer.Close()
	}
	if consumer != nil {
		consumer.Close()
	}
	if pulsarClient != nil {
		pulsarClient.Close()
	}
}

// Producer returns the Pulsar producer (nil when access == consume).
func Producer() pulsar.Producer {
	return producer
}

// Consumer returns the Pulsar consumer (nil when access == produce).
func Consumer() pulsar.Consumer {
	return consumer
}

// Receive blocks until a message is available (consume access only).
func Receive(ctx context.Context) (pulsar.Message, error) {
	if consumer == nil {
		return nil, fmt.Errorf("messaging: not configured for consume access")
	}
	return consumer.Receive(ctx)
}
