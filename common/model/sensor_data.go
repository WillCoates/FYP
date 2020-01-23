package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SensorData struct {
	ID        primitive.ObjectID `bson:"_id"`
	User      string             `bson:"user"`
	MessageID string             `bson:"messageid"`
	UnitID    string             `bson:"unitid"`
	Sensor    string             `bson:"sensor"`
	Value     string             `bson:"value"`
	Timestamp int64              `bson:"timestamp"`
}
