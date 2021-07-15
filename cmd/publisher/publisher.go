package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func quitOnFailure(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ": ", err)
	}
}

func main() {
	fmt.Println("Hello Publisher!")

	taskc := make(chan string, 10)

	go func() {
		r := bufio.NewReader(os.Stdin)
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			str = str[:len(str)-1]
			if str == "STOP" {
				break
			}
			taskc <- str
		}
		close(taskc)
	}()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	quitOnFailure(err, "Ошибка Dial")
	defer conn.Close()

	ch, err := conn.Channel()
	quitOnFailure(err, "Ошибка Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mySecondQueue",
		true,  //durable,
		false, //autoDelete,
		false, //exclusive,
		false, //noWait,
		nil,   //args,
	)
	quitOnFailure(err, "Faild to declare a queue")

	for task := range taskc {
		body := []byte(task)

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         body,
				ContentType:  "plain/text",
			},
		)
		quitOnFailure(err, "Faild to publishing")
	}

}
