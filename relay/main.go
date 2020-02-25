package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/pelletier/go-toml"
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

func main() {
	var config Config

	configFile := "relay.toml"

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	configTree, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalln("Failed to load config", err)
	}
	configTree.Unmarshal(&config)

	dbConn, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		log.Println("Failed to create MongoDB client")
		log.Fatalln(err)
	}

	dbctx, dbcancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer dbcancel()

	err = dbConn.Connect(dbctx)
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		log.Fatalln(err)
	}

	db := dbConn.Database("fyp")

	for {
		var amqpConnection *amqp.Connection
		for {
			amqpUrl, err := url.Parse(config.AmqpURL)
			if err != nil {
				log.Fatalln("Failed to decode AMQP url", err)
			}

			amqpUrl.User = url.UserPassword(config.AuthUsername, config.AuthPassword)

			amqpConnection, err = amqp.Dial(amqpUrl.String())
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

		for i := 0; i < config.Listeners; i++ {
			go listener(amqpConnection, db)
		}

		closeErrors := amqpConnection.NotifyClose(make(chan *amqp.Error))

		for err = range closeErrors {
			log.Println("Disconnected from RabbitMQ with error")
			log.Println(err)
		}

		log.Println("Reconnecting to RabbitMQ")
	}
}
