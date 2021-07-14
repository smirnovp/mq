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
	fmt.Println("Hello, Consumer!")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	quitOnFailure(err, "Failure to open Connection")
	defer conn.Close()

	ch, err := conn.Channel()
	quitOnFailure(err, "Failure get channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("myFirstQueue", false, false, false, false, nil)
	quitOnFailure(err, "Failure to declare queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	quitOnFailure(err, "Failure to get msgs channel")

	go func() {
		for msg := range msgs {
			fmt.Println(string(msg.Body))
		}
	}()

	done := make(chan struct{})
	<-done

}
