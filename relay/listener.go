package main

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func listener(connection *amqp.Connection, collection *mongo.Collection) {
	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("Failed to create channel")
		return
	}
	defer ch.Close()
	msgs, err := ch.Consume("sensors", "", false, false, false, false, nil)

	for msg := range msgs {
		var payload SensorPayload
		err := json.Unmarshal(msg.Body, &payload)
		// TODO: Retrieve data from topic
		if err != nil {
			fmt.Println("Failed to decode payload")
			fmt.Println(err)
			msg.Nack(false, false) // Discard
		} else {

			msg.Ack(false)
			msg.Nack(false, true) // Requeue
		}
	}
}
