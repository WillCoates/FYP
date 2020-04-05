package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) UpdateField(ctx context.Context, field *model.SensorSite, users []string) error {
	var err error

	fields := logic.db.Collection("sites")
	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return err
		}
	}

	query := make(bson.M)
	query["user"] = bson.M{"$in": userIds}
	query["name"] = field.Name

	res, err := fields.UpdateOne(ctx, query, bson.M{"$set": field})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNoSuchField
	}

	return nil
}
