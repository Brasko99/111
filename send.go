package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"main",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"main1", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	fmt.Println(string(q.Name))

	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type EmailData struct {
		Email   string `json:"email"`
		Content string `json:"content"`
		Text    string `json:"text"`
	}

	Message := EmailData{
		Email:   "born_to_fight_63@mail.ru",
		Content: "ContentContent",
		Text:    "TextTextText",
	}

	MrshlMessage, _ := json.Marshal(Message)
	fmt.Println(string(MrshlMessage))

	err = ch.PublishWithContext(ctx,
		"logs_topic", // exchange
		"email",      // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(MrshlMessage),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", string(MrshlMessage))
}
