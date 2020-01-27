package business

import (
	"context"
	"errors"

	"github.com/WillCoates/FYP/auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrInvalidAudience is returned by methods when an invalid audience is specified
var ErrInvalidAudience error = errors.New("Invalid audience")

// GetAudience retrieves an audience from the database
func (logic *Logic) GetAudience(ctx context.Context, name string) (*model.Audience, error) {
	audiencesCollection := logic.db.Collection("audiences")
	result := audiencesCollection.FindOne(ctx, bson.M{"name": name})
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidAudience
		}
		return nil, err
	}

	audience := new(model.Audience)

	err = result.Decode(audience)

	if err != nil {
		return nil, err
	}

	return audience, nil
}
