package rabbit

import (
	"bytes"
	"encoding/json"

	"github.com/alfarady/snowball"
	"github.com/alfarady/snowball/exception"
	"github.com/streadway/amqp"
)

func DialAmqpInstance(config snowball.Config) amqp.Channel {
	conn, err := amqp.Dial(config.ConnectionString)
	exception.PanicIfNeeded(err)

	ch, err := conn.Channel()
	exception.PanicIfNeeded(err)

	// Default Exchange
	exchangeName := "events"
	exchangeKind := "topic"

	if config.Exchange.Name != "" {
		exchangeName = config.Exchange.Name
	}

	if config.Exchange.Type != "" {
		exchangeKind = config.Exchange.Type
	}

	err = ch.ExchangeDeclare(
		exchangeName,               // name
		exchangeKind,               // kind
		config.Exchange.Durable,    // durable
		config.Exchange.AutoDelete, // auto delete
		config.Exchange.Internal,   // internal
		config.Exchange.NoWait,     // no wait
		config.Exchange.Arguments,  // amqp.Table
	)
	exception.PanicIfNeeded(err)

	return *ch
}

func DeserializeMessage(b []byte) (map[string]interface{}, error) {
	var msg map[string]interface{}
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

func SerializeMessage(msg interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}
