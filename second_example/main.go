//Simple RabbitMQ Producer
package main

import (
	"log"

	"github.com/streadway/amqp"
)

//Error helper function
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main() {
	//1. create a connection to the amqp protocol
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//2. Create channel for API utilization
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	//3. Declare a queue
	q, err := ch.QueueDeclare(
		"Hello", //name of queue
		false,   //durable
		false,   //Delete when used
		false,   //exclusive
		false,   //No-wait
		nil,     //arguement
	)
	failOnError(err, "Failed to declare queue")

	//4Publish message
	body := "Hello Rabbit folk"
	err = ch.Publish(
		"",     //exchange
		q.Name, //routing key
		false,  //mandatory
		false,  //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}
