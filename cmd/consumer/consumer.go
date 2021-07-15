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

	err = ch.ExchangeDeclare(
		"myFirstExchange", //name,
		"fanout",          // kind
		true,              // durable
		false,             //autoDelete
		false,             //internal
		false,             //noWait
		nil,               //args
	)
	quitOnFailure(err, "Failure to declare an Exchange")

	q, err := ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //autoDelete
		true,  //exclusive
		false, //noWait
		nil,   //args
	)
	quitOnFailure(err, "Failure to declare a Queue")

	err = ch.QueueBind(
		q.Name,            //name
		"",                //key
		"myFirstExchange", //exchange
		false,             //noWait
		nil,               //args
	)
	quitOnFailure(err, "Failure to set Qos")

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer
		false,  //autoAck
		false,  //exclusive
		false,  //nolocal
		false,  //noWait
		nil,    //args
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
