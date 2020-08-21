package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

//error handling
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to rabitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Faild to create rabbitMQ channel")
	defer ch.Close()

	//declare queue just in case it hasn't been by the producer
	q, err := ch.QueueDeclare(
		"task_queue", //queue name
		true,         //durable
		false,        //delete when used
		false,        //exclusive
		false,        //no-wait
		nil,          //arguements
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     //prefetch count
		0,     //prefetch size
		false, //global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, //queue Name
		"",     //consumer
		false,  //auto-ack
		false,  //exclusive
		false,  //no-local
		false,  //no-wait
		nil,    //args
	)
	failOnError(err, "Failed to register consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit, press CTRL+C ")
	<-forever
}
