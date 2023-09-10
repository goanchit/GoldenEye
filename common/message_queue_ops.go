package common

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s %s", msg, err)
	}
}

func SendMessageToQ(ctx context.Context, qName string, message string) (bool, string) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed To Connect to Rabbit MQ")
	defer conn.Close()

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	// Declaring Queue as durable so messages aren't lost in case of mq going down
	q, err := ch.QueueDeclare(
		qName, true, false, false, false, nil,
	)

	failOnError(err, "Failed to declare a queue")

	defer ch.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	err = ch.PublishWithContext(
		ctx, "", q.Name, false, false, amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // This will make the queue messages persistent. These wont be lost in case of server going down
		},
	)

	if err != nil {
		return false, "Failed to publish a message"
	}

	log.Printf(" [x] Sent %s\n", message)
	return true, ""
}
