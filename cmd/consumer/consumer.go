package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Hello Consumer!")

	if len(os.Args) < 2 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "<command>.n.m")
		os.Exit(1)
	}

	var err error
	exitIfError := func(msg string) {
		if err != nil {
			log.Fatalln(msg, err)
		}
	}

	conn, err := amqp.Dial("amqp://localhost:5672")
	exitIfError("Failed to Dial")
	defer conn.Close()

	ch, err := conn.Channel()
	exitIfError("Failed to get Channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"myForthExchange", // name
		"topic",           // kind
		true,              // durable
		false,             // autoDelete
		false,             // internal
		false,             // noWait
		nil,               // args
	)
	exitIfError("Failed to Declare Exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	exitIfError("Failed to Declare Queue")

	key := strings.Join(os.Args[1:], ".")

	err = ch.QueueBind(
		q.Name,            //name
		key,               // key
		"myForthExchange", // exchange
		false,             // noWait
		nil,               // args
	)
	exitIfError("Failed to Bind Queue")

	msgsc, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		true,   //autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	exitIfError("Failed to Consume")

	for msg := range msgsc {
		b := msg.Body
		bs := bytes.Split(b, []byte("."))

		var ans string

		if len(bs) < 3 {
			ans = "Error: too few elements"
		} else {
			a, erra := strconv.Atoi(string(bs[1]))
			b, errb := strconv.Atoi(string(bs[2]))
			if erra != nil || errb != nil {
				ans = "Error convert to integer"
			} else {
				switch string(bs[0]) {
				case "sum":
					ans = strconv.Itoa(a + b)
				case "mul":
					ans = strconv.Itoa(a * b)
				case "sub":
					ans = strconv.Itoa(a - b)
				}
			}
		}
		fmt.Println(string(b), "ans = ", ans)

		err = ch.Publish(
			"",          // exchange
			msg.ReplyTo, // key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				CorrelationId: msg.CorrelationId,
				ContentType:   "plain/text",
				Body:          []byte(ans),
			}, //msg
		)
		if err != nil {
			log.Println("Failed to Publish", err)
		}
	}

}
