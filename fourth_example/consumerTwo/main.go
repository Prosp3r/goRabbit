package main

import (
	"log"

	"github.com/streadway/amqp"
)

//log errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", err, msg)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connection to rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel to rabbitMQ")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   //exchange name to connect to
		"fanout", //Kind of exchange
		true,     //durable
		false,    //Auto-delete
		false,    //Internal
		false,    //no-wait
		nil,      //arguments
	)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		"",    //name
		false, //durable
		false, //auto-delete when unused
		true,  //exclusive
		false, //no-wait
		nil,   //arguements
	)
	failOnError(err, "Failed to declare Queue")

	err = ch.QueueBind(
		q.Name, //Queue Name
		"",     //Routing Key
		"logs", //exchange
		false,  //no-wait
		nil,    //arguements
	)
	failOnError(err, "Failed to bind queue")

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     //consumer
		true,   //auto-acknowledge
		false,  //exclusive
		false,  //no-local
		false,  //noWait
		nil,    //arguements
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
