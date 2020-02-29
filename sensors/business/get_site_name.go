package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) GetSiteName(users []string, id *primitive.ObjectID) (string, error) {
	sites := logic.db.Collection("sites")
	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		var err error
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return "", err
		}
	}

	var site model.SensorSite

	err := sites.FindOne(context.Background(), bson.M{"user": bson.M{"$in": userIds}, "_id": id}).Decode(&site)
	if err != nil {
		return "", err
	}
	return site.Name, nil
}
