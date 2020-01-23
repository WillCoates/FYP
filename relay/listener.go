package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func listener(connection *amqp.Connection, collection *mongo.Collection) {
	topicRegex, err := regexp.Compile("^([a-zA-Z0-9]+)\\.sensor(\\.([a-zA-Z0-9]+))?(\\.[a-zA-Z0-9]+)*$")
	if err != nil {
		fmt.Println("Regex error")
		fmt.Println(err)
		return
	}

	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("Failed to create channel")
		fmt.Println(err)
		return
	}
	defer ch.Close()
	msgs, err := ch.Consume("sensors", "", false, false, false, false, nil)

	for msg := range msgs {
		topicData := topicRegex.FindStringSubmatch(msg.RoutingKey)
		if len(topicData) == 0 {
			log.Printf("Recieved badly formatted topic %s", msg.RoutingKey)
			msg.Nack(false, false) // Discard
			continue
		}
		log.Printf("Message from user %s from relay %s (%s)", topicData[1], topicData[3], topicData[0])

		var payload SensorPayload
		err := json.Unmarshal(msg.Body, &payload)

		if err != nil {
			log.Println("Failed to decode payload")
			log.Println(err)
			msg.Nack(false, false) // Discard
		} else {
			data := payload.ToSensorData()
			data.User = topicData[1]
			_, err = collection.InsertOne(context.Background(), data)
			if err != nil {
				log.Println("Failed to store data")
				log.Println(err)
				msg.Nack(false, true) // Requeue
			} else {
				msg.Ack(false)
			}
		}
	}

	log.Println("Server stopped sending messages")
}
