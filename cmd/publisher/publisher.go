package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Hello publisher!")

	var err error

	exitIfError := func(msg string) {
		if err != nil {
			log.Fatalln(msg, err)
		}
	}

	commandc := make(chan []byte)

	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			b, err := r.ReadBytes('\n')
			if err != nil {
				log.Println("Ошибка чтения команды из Stdin:", err)
				continue
			}
			if len(b) > 0 {
				b = b[:len(b)-1]
			}
			if string(b) == "STOP" {
				break
			}
			commandc <- b
		}
		close(commandc)
	}()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	exitIfError("Failed to Dial")
	defer conn.Close()

	ch, err := conn.Channel()
	exitIfError("Faild to get Channel")
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
	exitIfError("Faild to Daclare Exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)

	exitIfError(`Faild to Declare Queue "rpc_result"`)

	msgc, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	exitIfError("Failed to Consume")

	corrID := RandString(32)

	go func() {

		for comm := range commandc {

			c := strings.Replace(string(comm), " ", ".", -1)

			err := ch.Publish(
				"myForthExchange", // exchange
				c,                 // key
				false,             // mandatory
				false,             // immediate
				amqp.Publishing{
					ReplyTo:       q.Name,
					CorrelationId: corrID,
					ContentType:   "plain/text",
					Body:          []byte(c),
				}, // msg
			)
			if err != nil {
				log.Println("Failed to Publish", err)
			}
		}

		ch.Close()
	}()

	for msg := range msgc {
		if _, err := os.Stdout.Write(append(msg.Body, '\n')); err != nil {
			log.Println("Ошибка записи в Stdout:", err)
		}
	}

}

// RandString генерирует случайную строку длиной n состоящую только из букв
func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // максимум 63 знака
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	var src = rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			err := sb.WriteByte(letterBytes[idx])
			if err != nil {
				fmt.Println("RandString(): sb.WriteByte(): ", err)
			}
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
