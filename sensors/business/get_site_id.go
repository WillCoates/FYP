package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (logic *Logic) GetSiteId(users []string, name string, create bool) (primitive.ObjectID, error) {
	// Allow sensors without sites
	if name == "" {
		return primitive.NilObjectID, nil
	}

	sites := logic.db.Collection("sites")
	userIds := make([]primitive.ObjectID, len(users))

	for i, user := range users {
		var err error
		userIds[i], err = primitive.ObjectIDFromHex(user)
		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	var site model.SensorSite

	err := sites.FindOne(context.Background(), bson.M{"user": bson.M{"$in": userIds}, "name": name}).Decode(&site)
	if err == mongo.ErrNoDocuments && create {
		res, err := sites.InsertOne(context.Background(), bson.M{"user": userIds[0], "name": name})
		if err != nil {
			return primitive.NilObjectID, err
		} else {
			return res.InsertedID.(primitive.ObjectID), nil
		}
	} else if err != nil {
		return primitive.NilObjectID, err
	}
	return site.ID, nil
}
