package business

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/bson"
)

// InvalidateToken invalidates a token within the database
func (logic *Logic) InvalidateToken(ctx context.Context, token *auth.Token) error {
	tokensCollection := logic.db.Collection("tokens")

	_, err := tokensCollection.UpdateOne(ctx, bson.M{"jwtid": token.Payload.JwtID}, bson.M{"$set": bson.M{"valid": false}})

	if err != nil {
		return err
	}

	return nil
}
