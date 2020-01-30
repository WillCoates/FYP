package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthCode struct {
	ID        *primitive.ObjectID `bson:"_id,omitempty"`
	Code      string              `bson:"code"`
	Token     string              `bson:"token"`
	Expires   int64               `bson:"expires"`
	Redirect  string              `bson:"redirect"`
	Claimed   bool                `bson:"claimed"`
	Challenge *Challenge          `bson:"challenge"`
}
