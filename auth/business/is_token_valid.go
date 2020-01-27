package business

import (
	"context"
	"errors"
	"time"

	"github.com/WillCoates/FYP/auth/model"
	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrTokenInvalid is returned by IsTokenValid when the token has been
// invalidated or expired
var ErrTokenInvalid error = errors.New("Token invalid")

// IsTokenValid checks the database to ensure the token is still valid,
// returning ErrTokenInvalid if it isn't
func (logic *Logic) IsTokenValid(ctx context.Context, token *auth.Token) error {
	// Expiry
	now := time.Now().Unix()
	if token.Payload.Expires >= now || token.Payload.NotBefore < now {
		return ErrTokenInvalid
	}

	// Invalidation
	tokensCollection := logic.db.Collection("tokens")

	var registeredToken model.Token
	err := tokensCollection.FindOne(ctx, bson.M{"jwtid": token.Payload.JwtID}).Decode(&registeredToken)

	if err == mongo.ErrNoDocuments {
		return ErrTokenInvalid
	} else if err != nil {
		return err
	}

	if !registeredToken.Valid {
		return ErrTokenInvalid
	}

	return nil
}
