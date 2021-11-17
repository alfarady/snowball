# Snowball

<p align="center"><img src="doc/snowball-logo.png" width="360"></p>

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
  + [Importing the package](#importing-the-package)
  + [Making a simple `PUBLISHER` message](#making-a-simple-publisher-message)
  + [Making a simple `CONSUMER` message](#making-a-simple-consumer-message)
* [License](#license)

## Description

Snowball is an RabbitMQ event client that helps your application publish or consume message using RabbitMQ. All HTTP methods are exposed as a fluent interface.

## Installation
```
go get -u github.com/alfarady/snowball
```

## Usage

### Importing the package

This package can be used by adding the following import statement to your `.go` files.

```go
import "github.com/alfarady/snowball"
import "github.com/alfarady/snowball/rabbit"
```

### Making a simple `PUBLISHER` message
The below example will publish message using routing key:

```go
// Create Snowball Config
config := snowball.Config{
    ConnectionString: "amqp://rabbitmq:rabbitmq@localhost:5672/",
    MaxRetry:         1,
}

// Create Snowball Client
snowballClient := rabbit.NewSnowball(config)

// Do publish message, "test.something" is a routing key
snowballClient.Publish("test.something", snowball.Publishable{
    Body: Model{
        Test: 1234,
    },
})
```

### Making a simple `CONSUMER` message
The below example will consume message by routing key:

```go
// Your main function
func main() {
	config := snowball.Config{
		ConnectionString: "amqp://rabbitmq:rabbitmq@localhost:5672/",
		MaxRetry:         1,
	}

	snowballClient := rabbit.NewSnowball(config)

	snowballClient.Publish("test.something", snowball.Publishable{
		Body: Model{
			Test: 1234,
		},
	})

    // You should make your own channel
	forever := make(chan bool)

    // Call using go routines
	go func() {
		snowballClient.Listen("test.something", Handle)
	}()

    fmt.Print("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Handler function, this func will called when message coming
func Handle(message snowball.Listenable) error {
	fmt.Print(message.Body)

    // Return nil means there is no need to retry, you can return err to retry the queue
	return nil
}
```

## License

```
Copyright 2021, Alfarady (https://alfais.online)
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```