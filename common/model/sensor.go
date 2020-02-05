package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Sensor struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	User   primitive.ObjectID `bson:"user,omitempty"`
	UnitID string             `bson:"unitid"`
	Name   string             `bson:"name"`
	Info   *SensorInfo        `bson:"info"`
}
