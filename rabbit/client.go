package rabbit

import (
	"fmt"
	"strings"

	"github.com/alfarady/snowball"
	"github.com/alfarady/snowball/exception"
	"github.com/streadway/amqp"
)

const DefaultExchangeName = "events"
const DefaultContentType = "text/plain"

type Client struct {
	Config  snowball.Config
	Channel amqp.Channel
}

func NewSnowball(config snowball.Config) snowball.Client {
	client := &Client{
		Config:  config,
		Channel: amqp.Channel{},
	}

	client = client.init()

	return client
}

func (client *Client) init() *Client {
	client.Channel = DialAmqpInstance(client.Config)

	return client
}

func (client *Client) Publish(routingKey string, publishable snowball.Publishable) error {
	defaultExchangeName := DefaultExchangeName
	defaultAttempts := 1
	defaultHeaders := make(map[string]interface{})

	if publishable.ExchangeName != "" {
		defaultExchangeName = publishable.ExchangeName
	}

	if publishable.Headers != nil {
		defaultHeaders = publishable.Headers
	}

	body, err := SerializeMessage(publishable.Body)

	if err != nil {
		return err
	}

	err = client.doPublish(
		defaultExchangeName,
		routingKey,
		DefaultContentType,
		body,
		defaultHeaders,
		int32(defaultAttempts),
	)

	return err
}

func (client *Client) doPublish(
	exchangeName string,
	routingKey string,
	contentType string,
	body []byte,
	headers map[string]interface{},
	attempts int32,
) error {
	headers["x-attempts"] = attempts

	err := client.Channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
			Headers:     headers,
		},
	)

	return err
}

func (client *Client) Listen(routingKey string, handler snowball.Handler, exchangeName ...string) error {
	defaultQueuePrefix := "Snowball"
	defaultExchangeName := DefaultExchangeName

	if client.Config.Queue.Prefix != "" {
		defaultQueuePrefix = client.Config.Queue.Prefix
	}

	if len(exchangeName) > 0 {
		defaultExchangeName = exchangeName[0]
	}

	q, err := client.Channel.QueueDeclare(
		fmt.Sprintf("%s:%s", defaultQueuePrefix, routingKey),
		client.Config.Exchange.Durable,
		client.Config.Exchange.AutoDelete,
		false,
		client.Config.Exchange.NoWait,
		client.Config.Exchange.Arguments,
	)
	exception.PanicIfNeeded(err)

	err = client.Channel.QueueBind(
		q.Name,
		routingKey,
		defaultExchangeName,
		client.Config.Exchange.NoWait,
		client.Config.Exchange.Arguments,
	)
	exception.PanicIfNeeded(err)

	messages, err := client.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	exception.PanicIfNeeded(err)

	go func() {
		for message := range messages {
			fmt.Printf("[*] Received message: %s - %s", routingKey, strings.ReplaceAll(string(message.Body), "\n", ""))
			attempts := message.Headers["x-attempts"].(int32)

			body, _ := DeserializeMessage(message.Body)
			payload := snowball.Listenable{
				Body:    body,
				Headers: message.Headers,
			}

			err = handler(payload)

			if err != nil {
				if attempts < client.Config.MaxRetry {
					client.doPublish(
						defaultExchangeName,
						routingKey,
						message.ContentType,
						message.Body,
						message.Headers,
						attempts+1,
					)
				} else {
					fmt.Printf("Too many attempts: %s - %s", routingKey, message.Body)
				}
			}
		}
	}()

	return nil
}
