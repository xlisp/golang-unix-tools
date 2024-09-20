package main

import (
	"testing"
	"github.com/streadway/amqp"
)

func TestSendMessage(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to declare a queue: %s", err)
	}

	messages := make(chan string)

	go func() {
		for msg := range messages {
			body := msg
			err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)
			if err != nil {
				t.Fatalf("Failed to publish a message: %s", err)
			}
		}
	}()

	messages <- "Hello World!"
	close(messages)
}