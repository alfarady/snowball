package rabbit

import (
	"testing"

	"github.com/alfarady/snowball"
	"github.com/stretchr/testify/assert"
)

func TestSerializeMessage(t *testing.T) {
	assert := assert.New(t)

	message := map[string]interface{}{
		"test": 1234,
	}
	stringMessage := `{"test":1234}`
	serializedMessage, err := SerializeMessage(message)

	assert.Nil(err)
	assert.NotNil(serializedMessage)
	assert.JSONEq(stringMessage, string(serializedMessage))

}

func TestDeserializeMessage(t *testing.T) {
	assert := assert.New(t)

	message := map[string]interface{}{
		"test": float64(1234),
	}
	stringMessage := `{"test":1234}`
	deserializedMessage, err := DeserializeMessage([]byte(stringMessage))

	assert.Nil(err)
	assert.NotNil(deserializedMessage)
	assert.Equal(message["test"], deserializedMessage["test"].(float64))

}

func TestDialAmqpInstance(t *testing.T) {
	assert := assert.New(t)

	config := snowball.Config{
		ConnectionString: "amqp://rabbitmq:rabbitmq@localhost:5672/",
	}

	client := DialAmqpInstance(config)

	assert.NotNil(client)
}
