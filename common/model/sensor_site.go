package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SensorSite struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	User primitive.ObjectID `bson:"user,omitempty"`
	Name string             `bson:"name"`
}
