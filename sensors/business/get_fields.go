package business

import (
	"context"
	"log"

	"github.com/WillCoates/FYP/common/model"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) GetFields(ctx context.Context, users []string, fieldQuery bson.M) (chan model.SensorSite, error) {
	var fields = logic.db.Collection("sites")
	var err error

	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return nil, err
		}
	}

	fieldQuery = util.CloneMapStringIface(fieldQuery)
	fieldQuery["user"] = bson.M{"$in": userIds}

	result, err := fields.Find(ctx, fieldQuery)
	if err != nil {
		return nil, err
	}

	fieldsChan := make(chan model.SensorSite)

	go func() {
		defer result.Close(ctx)
		for result.Next(ctx) {
			var field model.SensorSite
			err := result.Decode(&field)
			if err != nil {
				log.Println("Failed to decode field")
				log.Println(err)
			} else {
				fieldsChan <- field
			}
		}
		close(fieldsChan)
	}()

	return fieldsChan, nil
}
