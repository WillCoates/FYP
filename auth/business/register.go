package business

import (
	"context"
	"errors"

	"github.com/WillCoates/FYP/auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrAlreadyRegistered error = errors.New("Email address already in use")

func (logic *Logic) Register(ctx context.Context, email, name, password string) error {
	usersCollection := logic.db.Collection("users")

	count, err := usersCollection.CountDocuments(ctx, bson.M{"email": email})

	if err != nil {
		return err
	}
	if count != 0 {
		return ErrAlreadyRegistered
	}

	var user model.User
	user.ID = primitive.NewObjectID()
	user.EmailAddress = email
	user.Name = name
	user.SetPassword(password)

	_, err = usersCollection.InsertOne(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
