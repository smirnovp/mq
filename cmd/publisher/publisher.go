package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func quitOnFailure(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ": ", err)
	}
}

func main() {
	fmt.Println("Hello Publisher!")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	quitOnFailure(err, "Ошибка Dial")
	defer conn.Close()

	ch, err := conn.Channel()
	quitOnFailure(err, "Ошибка Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"myFirstQueue",
		false, //durable,
		false, //autoDelete,
		false, //exclusive,
		false, //noWait,
		nil,   //args,
	)
	quitOnFailure(err, "Faild to declare a queue")

	body := []byte("Нелло пиплы!")

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			ContentType: "plain/text",
		},
	)
	quitOnFailure(err, "Faild to publishing")
}
