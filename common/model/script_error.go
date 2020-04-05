package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ScriptError struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Script    primitive.ObjectID `bson:"script,omitempty"`
	Message   string             `bson:"message"`
	Timestamp int64              `bson:"timestamp"`
}
