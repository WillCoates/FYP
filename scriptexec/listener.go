package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func listener(connection *amqp.Connection, db *mongo.Database) {
	ch, err := connection.Channel()
	if err != nil {
		fmt.Println("Failed to create channel")
		fmt.Println(err)
		return
	}
	defer ch.Close()
	msgs, err := ch.Consume("sensorscripting", "", false, false, false, false, nil)

	for msg := range msgs {
		log.Printf("Message %s", msg.RoutingKey)
		firstDotOffset := strings.Index(msg.RoutingKey, ".")
		var userID, topic string
		if firstDotOffset == -1 {
			userID = msg.RoutingKey
			topic = ""
		} else {
			userID = msg.RoutingKey[:firstDotOffset]
			topic = msg.RoutingKey[firstDotOffset+1:]
		}

		log.Printf("Message from user %s (%s)", userID, topic)

		if err != nil {
			log.Println("Failed to decode payload", err)
			msg.Nack(false, false)
			continue
		}

		uid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			log.Println("Failed to decode user ID", err)
			msg.Nack(false, false)
			continue
		}

		err, requeue := ExecScript(db, ch, uid, topic, string(msg.Body))
		if err != nil {
			log.Println("Failed to execute script", err)
			msg.Nack(false, requeue)
			continue
		}
		msg.Ack(false)
	}
}
