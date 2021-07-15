package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		"mySecondExchange", //name,
		"direct",           // kind
		true,               // durable
		false,              //autoDelete
		false,              //internal
		false,              //noWait
		nil,                //args
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

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [a] [b] [c] [d]\n", filepath.Base(os.Args[0]))
	}
	keys := os.Args[1:]

	for _, key := range keys {
		err = ch.QueueBind(
			q.Name,             //name
			key,                //key
			"mySecondExchange", //exchange
			false,              //noWait
			nil,                //args
		)
		quitOnFailure(err, "Failure to set Qos")
	}

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
			//t := time.Duration(l)
			//fmt.Println(string(b), ":", l, "sec.")
			fmt.Println(string(b), "len=", l)
			//time.Sleep(t * time.Second)
			err := msg.Ack(false)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	done := make(chan struct{})
	<-done

}
