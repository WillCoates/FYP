package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}

func prepareMessageQueue(connection *amqp.Connection) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Persistant queue
	_, err = ch.QueueDeclare("sensors", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.QueueBind("sensors", "*.sensor.#", "amq.topic", false, nil)
	return err
}

/*
func main() {
	// TODO: Load configuration
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://relay:Y&q&tdPuX2G1_4G8@docker:27017"))
	checkError(err, "Failed to create MongoDB client")

	dbctx, dbcancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer dbcancel()

	err = db.Connect(dbctx)
	checkError(err, "Failed to connect to MongoDB")

	collection := db.Database("fyp").Collection("sensor_readings")

	conn, err := amqp.Dial("amqp://docker:5672/")
	checkError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	err = prepareMessageQueue(conn)
	checkError(err, "Failed to create message queue")

	halt := make(chan bool)

	for i := 0; i < 1; i++ {
		go listener(conn, collection)
	}

	<-halt
}
*/

func main() {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://relay:Y&q&tdPuX2G1_4G8@docker:27017"))
	if err != nil {
		log.Println("Failed to create MongoDB client")
		log.Fatalln(err)
	}

	dbctx, dbcancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer dbcancel()

	err = db.Connect(dbctx)
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		log.Fatalln(err)
	}

	collection := db.Database("fyp").Collection("sensor_readings")

	for {
		var amqpConnection *amqp.Connection
		for {
			amqpConnection, err = amqp.Dial("amqp://docker:5672/")
			if err != nil {
				log.Println("Failed to connect to RabbitMQ")
				log.Println(err)
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}

		log.Println("Connected to RabbitMQ")

		err = prepareMessageQueue(amqpConnection)
		if err != nil {
			log.Println("Failed to register message queue")
			log.Fatalln(err)
		}

		for i := 0; i < 1; i++ {
			go listener(amqpConnection, collection)
		}

		closeErrors := amqpConnection.NotifyClose(make(chan *amqp.Error))

		for err = range closeErrors {
			log.Println("Disconnected from RabbitMQ with error")
			log.Println(err)
		}

		log.Println("Reconnecting to RabbitMQ")
	}
}
