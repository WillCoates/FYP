package model

import (
	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID      primitive.ObjectID `bson:"_id"`
	JwtID   string             `bson:"jwtid"`
	Subject primitive.ObjectID `bson:"subject"`
	Expires int64              `bson:"expires"`
	Valid   bool               `bson:"valid"`
}

func CreateRegisteredToken(token *auth.Token) (Token, error) {
	var registeredToken Token
	var err error

	registeredToken.ID = primitive.NewObjectID()
	registeredToken.JwtID = token.Payload.JwtID
	registeredToken.Subject, err = primitive.ObjectIDFromHex(token.Payload.Subject)
	registeredToken.Expires = token.Payload.Expires
	registeredToken.Valid = true

	return registeredToken, err
}
