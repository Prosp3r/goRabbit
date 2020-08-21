//Demonstrate Work Queues in rabbitMQ
package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Could not connect to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", //name
		true,         //durable
		false,        //delete when used
		false,        //exclusive
		false,        //no-wait
		nil,          //arguements
	)
	failOnError(err, "Failed to declare queue")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  //mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish message")
	log.Printf(" [x] Sent %s", body)
}

//
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "yxHello..............."
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
