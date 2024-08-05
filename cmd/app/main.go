package main

import (
	"github.com/rabbitmq/amqp091-go"
	"os"
)

func main() {
	url := os.Getenv("AMQP_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	connection, err := amqp091.Dial(url)
}
