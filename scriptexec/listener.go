package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func listener(connection *amqp.Connection, db *mongo.Database) {
	topicRegex := regexp.MustCompile("^([a-zA-Z0-9]+)\\.sensor(\\.([a-zA-Z0-9]+))?(\\.[a-zA-Z0-9]+)*$")

	_ = db.Collection("scripts")

	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("Failed to create channel")
		fmt.Println(err)
		return
	}
	defer ch.Close()
	msgs, err := ch.Consume("sensorscripting", "", false, false, false, false, nil)

	for msg := range msgs {
		topicData := topicRegex.FindStringSubmatch(msg.RoutingKey)
		if len(topicData) == 0 {
			log.Printf("Recieved badly formatted topic %s", msg.RoutingKey)
			msg.Nack(false, false) // Discard
			continue
		}
		log.Printf("Message from user %s from relay %s (%s)", topicData[1], topicData[3], topicData[0])

		data := make(map[string]interface{})
		err := json.Unmarshal(msg.Body, &data)

		if err != nil {
			log.Println("Failed to decode payload", err)
			msg.Nack(false, false)
			continue
		}
		_, err = primitive.ObjectIDFromHex(topicData[1])

		if err != nil {
			log.Println("Failed to decode user ID", err)
			continue
		}
	}
}
