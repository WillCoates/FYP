package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) UpdateSensor(ctx context.Context, sensor *model.Sensor, users []string) error {
	var err error

	sensors := logic.db.Collection("sensors")
	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return err
		}
	}

	query := make(bson.M)
	query["user"] = bson.M{"$in": userIds}
	query["unitid"] = sensor.UnitID
	query["info.sensor"] = sensor.Info.Sensor

	res, err := sensors.UpdateOne(ctx, query, bson.M{"$set": sensor})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrNoSuchSensor
	}

	return nil
}
