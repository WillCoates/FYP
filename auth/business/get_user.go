package business

import (
	"context"

	"github.com/WillCoates/FYP/auth/model"
	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) GetUser(ctx context.Context, token *auth.Token) (*model.User, error) {
	users := logic.db.Collection("users")
	user := new(model.User)

	id, err := primitive.ObjectIDFromHex(token.Payload.Subject)
	if err != nil {
		return nil, err
	}

	err = users.FindOne(ctx, bson.M{"_id": id}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
