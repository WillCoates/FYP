package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/WillCoates/FYP/common/model"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func listener(connection *amqp.Connection, db *mongo.Database) {
	topicRegex, err := regexp.Compile("^([a-zA-Z0-9]+)\\.sensor(\\.([a-zA-Z0-9]+))?(\\.[a-zA-Z0-9]+)*$")
	if err != nil {
		fmt.Println("Regex error")
		fmt.Println(err)
		return
	}

	sensors := db.Collection("sensors")
	readings := db.Collection("sensor_readings")
	types := db.Collection("sensor_types")

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
			var sensor model.Sensor
			var data model.SensorData
			var sensorID primitive.ObjectID

			userID, err := primitive.ObjectIDFromHex(topicData[1])

			if err != nil {
				log.Println("Failed to decode user ID")
				log.Println(err)
				msg.Nack(false, false) // Discard
			}

			err = sensors.FindOne(context.Background(), bson.M{"user": userID, "unitid": payload.UnitID, "info.sensor": payload.Sensor}).Decode(&sensor)
			if err == mongo.ErrNoDocuments {
				var info model.SensorInfo
				sensor.User = userID
				sensor.UnitID = payload.UnitID
				sensor.Name = payload.UnitID

				err = types.FindOne(context.Background(), bson.M{"sensor": payload.Sensor}).Decode(&info)
				if err == mongo.ErrNoDocuments {
					info.Sensor = payload.Sensor
					info.Measurement = payload.Sensor
					info.Units = payload.Sensor
					info.Hidden = false
				} else if err != nil {
					log.Println("Failed to fetch zero-config sensor info")
					log.Println(err)
					msg.Nack(false, true)
					continue
				}

				info.ID = primitive.ObjectID{}
				sensor.Info = &info

				res, err := sensors.InsertOne(context.Background(), &sensor)
				if err != nil {
					log.Println("Failed to create sensor")
					log.Println(err)
					msg.Nack(false, true)
					continue
				}
				sensorID = res.InsertedID.(primitive.ObjectID)
			} else if err != nil {
				log.Println("Failed to fetch sensor")
				log.Println(err)
				msg.Nack(false, true)
				continue
			} else {
				sensorID = sensor.ID
			}

			result := readings.FindOne(context.Background(), bson.M{"msgid": payload.MessageID, "sensor": sensorID})
			err = result.Err()
			if err == mongo.ErrNoDocuments {
				data.Sensor = sensorID
				data.Timestamp = payload.Timestamp
				data.Value = payload.Value
				data.MessageID = payload.MessageID

				_, err = readings.InsertOne(context.Background(), data)
				if err != nil {
					log.Println("Failed to store data")
					log.Println(err)
					msg.Nack(false, true) // Requeue
				} else {
					msg.Ack(false)
				}
			} else if err != nil {
				log.Println("Failed to check for duplicate reading", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}

	log.Println("Server stopped sending messages")
}
