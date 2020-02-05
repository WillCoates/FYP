package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SensorData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Sensor    primitive.ObjectID `bson:"sensor,omitempty"`
	MessageID string             `bson:"msgid"`
	Value     string             `bson:"value"`
	Timestamp int64              `bson:"timestamp"`
}
