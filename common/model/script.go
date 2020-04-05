package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Script struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	User         primitive.ObjectID `bson:"user,omitempty"`
	Name         string             `bson:"name"`
	Patterns     []string           `bson:"patterns,omitempty"`
	Source       string             `bson:"source"`
	LastModified int64              `bson:"lastmodified"`
}
