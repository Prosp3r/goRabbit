package main

import (
	"os"

	"github.com/streadway/amqp"
)

func main() {
	//get url for amqp connection
	url := os.Getenv("AMQP_URL")

	//Check if environment variable does not exist, user the default url
	if url == "" {
		//not to be done in production
		url = "amqp://guest:gues@localhost:5672"
	}

	//connect to rabbitMQ
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("Could not establish connection with rabbitMQ: " + err.Error())
	}


	//Create a channel from the Connection--Channels share a single TCP connection
	chan, err := connection.Channel()
	if err != nil {
		panic("Error opening rabbitMQ channel:" + err.Error())
	}
	//Declare an exchange topic event
	err = chan.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	//Declare a queue and bind it to the events exchange

	//create message ...must an instance of an amqp struct
	message := amqp.Publishing{
		Body: []byte("Hello Rabbit World"),
	}

	//publish message to exchange
	err = chan.Publish("events", "random-key", false, false, message)
	if err != nil{
		panic("error publishing message to the queue:" + err.Error())
	}

	//Declare/Create a Queue
	_, err := chan.QueueDeclare("test", true, false, false, false, nil)
	if err != nil {
		panic("Error delcaring Queue: "+ err.Error())
	}

	//Bind Queue to exhcnage 
	err = chan.QueueBind("test", "#", "events", false, nil)
	if err != nil {
		panic("Error binding queue to exchange : " + err.Error())
	}


}
