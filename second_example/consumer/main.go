package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	//Make a connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Create a channel connection
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a rabbitMQ channel")
	defer ch.Close()

	//Queue is declared here as well because it ensures the queue exists before consuming it.
	//THe producer may have not declared it at the time of connecting
	q, err := ch.QueueDeclare(
		"Hello", //Name
		false,   //durable
		false,   //Auto delete when used
		false,   //exclusive
		false,   //No-wait
		nil,     //Arguements
	)
	failOnError(err, "Failed to declare a queue")

	//
	msgs, err := ch.Consume(
		q.Name, //Name
		"",     //Consumer string
		true,   //Auto Ack ...acknowledge
		false,  //exclusive
		false,  //NoLocal
		false,  //no-wait consume immediately
		nil,    //args
	)
	failOnError(err, "Failed to register a consumer")

	//keepAlive channel
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit, press CTRL+C")
	<-forever
}
