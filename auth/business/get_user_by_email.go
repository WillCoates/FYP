package business

import (
	"context"

	"github.com/WillCoates/FYP/auth/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (logic *Logic) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	users := logic.db.Collection("users")
	user := new(model.User)

	err := users.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
