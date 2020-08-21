package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

//Log errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err)
	}
}

func main() {
	//Establish connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to rabbitMQ")
	defer conn.Close()

	//Channel
	ch, err := conn.Channel()
	failOnError(err, "Could not create channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   //Name of exchange
		"fanout", //type of exchange
		true,     //durable
		false,    //auto-delete
		false,    //internal
		false,    //no-wait
		nil,      //arguments
	)
	failOnError(err, "Failed to declare exchange")

	body := bodyFrom(os.Args)

	err = ch.Publish(
		"logs", //Exchange
		"",     //routing key
		false,  //Mandatory
		false,  //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish message")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "Hello rabbits triple time..."
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
