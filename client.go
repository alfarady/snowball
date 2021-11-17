package snowball

import "github.com/streadway/amqp"

type Config struct {
	ConnectionString string
	Exchange         Exchange
	Queue            Queue

	MaxRetry int32
}

type Exchange struct {
	Name       string
	Type       string
	Passive    bool
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Arguments  amqp.Table
}

type Queue struct {
	Prefix string
}

type Publishable struct {
	ExchangeName string
	Body         interface{}
	Headers      map[string]interface{}
}

type Listenable struct {
	Body    map[string]interface{}
	Headers map[string]interface{}
}

type Handler = func(payload Listenable) error

type Client interface {
	Publish(routingKey string, publishable Publishable) error
	Listen(routingKey string, handler Handler, exchangeName ...string) error
}
