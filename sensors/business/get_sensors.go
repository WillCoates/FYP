package business

import (
	"context"
	"log"

	"github.com/WillCoates/FYP/common/model"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) GetSensors(ctx context.Context, query bson.M, users ...string) (chan model.Sensor, error) {
	var err error

	sensors := logic.db.Collection("sensors")
	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return nil, err
		}
	}

	query = util.CloneMapStringIface(query)
	query["user"] = bson.M{"$in": userIds}

	result, err := sensors.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	channel := make(chan model.Sensor)

	go func() {
		defer result.Close(ctx)
		for result.Next(ctx) {
			var sensor model.Sensor
			err := result.Decode(&sensor)

			if err != nil {
				log.Println("Error decoding sensor")
				log.Println(err)
			} else {
				channel <- sensor
			}
		}
		close(channel)
	}()

	return channel, nil
}
