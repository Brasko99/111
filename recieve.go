package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Функция выдачи ошибки
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//создаём connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//создаём канал
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//создаём Exchange
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

	//создаём очередь Queue
	q, err := ch.QueueDeclare(
		"mail_service", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//делаем Binding Queue
	err = ch.QueueBind(
		q.Name,    // queue name
		"email.#", // routing key
		"main",    // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	//Получаем сообщение
	msgs, err := ch.Consume(
		q.Name, // queue
		"mail", // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	//обрабатываем сообщение
	go func() {
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)
			//объявляем структуру для принятого сообщения
			type EmailData struct {
				Email   string `json:"email"`
				Content string `json:"content"`
				Text    string `json:"text"`
			}
			//принимаем сообщение в переменную
			ReceivedMessage := d.Body

			//делаем Unmarshall сообщения
			var ReceivedStruct EmailData
			if err := json.Unmarshal(ReceivedMessage, &ReceivedStruct); err != nil {
				panic(err)
			}

			//создание переменных для передачи в html
			var Email string = ReceivedStruct.Email
			var Content string = ReceivedStruct.Content
			var Text string = ReceivedStruct.Text

			//Email получателя должен быть массивом по правилам использования функции отправки сообщения
			Email1 := []string{Email}

			// Данные отправителя сообщения
			from := "Olegovich99@inbox.ru"
			password := "58PGx1zk5Zxgb4WPwq5i"

			// Настройки smtp server configuration.
			smtpHost := "smtp.mail.ru"
			smtpPort := "587"

			// Аутентификация в почте отправителем.
			auth := smtp.PlainAuth("", from, password, smtpHost)

			// Передача переменных в html
			t, _ := template.ParseFiles("template.html")
			var body bytes.Buffer
			mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
			body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

			t.Execute(&body, struct {
				Email   string
				Content string
				Text    string
			}{
				Email,
				Content,
				Text,
			})

			// Отправка сообщения
			err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, Email1, body.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Email Sent!")

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
