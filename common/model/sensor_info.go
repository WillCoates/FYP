package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SensorInfo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Sensor      string             `bson:"sensor"`
	Measurement string             `bson:"measurement"`
	Units       string             `bson:"units"`
	Hidden      bool               `bson:"hidden"`
}
