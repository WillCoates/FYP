package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) DeleteField(ctx context.Context, field *model.SensorSite, users []string) error {
	var err error

	sensors := logic.db.Collection("sensors")
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

	var deletedField model.SensorSite

	err = fields.FindOneAndDelete(ctx, query).Decode(&deletedField)
	if err != nil {
		return err
	}

	sensors.UpdateMany(ctx, bson.M{"site": deletedField.ID}, bson.M{"$unset": bson.M{"site": ""}})

	return nil
}
