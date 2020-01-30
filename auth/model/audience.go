package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Audience struct {
	ID              primitive.ObjectID `bson:"_id"`
	Internal        bool               `bson:"internal"`
	Name            string             `bson:"name"`
	Perms           []string           `bson:"perms"`
	DefaultDuration uint32             `bson:"defaultExpires"`
	MaxDuration     uint32             `bson:"maxExpires"`
	Comment         string             `bson:"comment"`
	Secret          string             `bson:"secret"`
}
