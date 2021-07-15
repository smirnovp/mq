package main

import (
	"fmt"
	"log"
	"time"

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

	q, err := ch.QueueDeclare(
		"mySecondQueue",
		true,
		false,
		false,
		false,
		nil,
	)
	quitOnFailure(err, "Failure to declare queue")

	err = ch.Qos(
		1,
		0,
		false,
	)
	quitOnFailure(err, "Failure to set Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	quitOnFailure(err, "Failure to get msgs channel")

	go func() {
		for msg := range msgs {
			b := msg.Body
			l := len(b)
			t := time.Duration(l)
			fmt.Println(string(b), ":", l, "sec.")
			time.Sleep(t * time.Second)
			err := msg.Ack(false)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	done := make(chan struct{})
	<-done

}
