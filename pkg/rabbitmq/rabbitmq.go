package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type Option func(*config)

type config struct {
}

type Rabbit struct {
	*amqp091.Connection
}

func New(connString string, opts ...Option) (*Rabbit, error) {
	connection, err := amqp091.Dial(connString)
	if err != nil {
		return nil, err
	}

	return &Rabbit{Connection: connection}, nil
}
