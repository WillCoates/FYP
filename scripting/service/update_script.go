package service

import (
	"context"
	"time"

	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service *ScriptingService) UpdateScript(ctx context.Context, req *proto.Script) (*proto.Script, error) {
	scripts := service.db.Collection("scripts")

	_, perms, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	var users []string

	for _, perm := range perms {
		if perm.Permission == "manageScripting" {
			users = append(users, perm.For)
		}
	}

	if len(users) == 0 {
		return nil, ErrNoPermission
	}

	userIDs, err := util.StringIDToObjectID(users)
	if err != nil {
		return nil, err
	}

	scriptID, err := primitive.ObjectIDFromHex(req.Details.Id)
	if err != nil {
		return nil, err
	}

	lastMod := time.Now().Unix()

	filter := bson.M{
		"_id":  scriptID,
		"user": bson.M{"$in": userIDs},
	}

	update := bson.M{
		"name":         "",
		"patterns":     "",
		"source":       "",
		"lastmodified": lastMod,
	}

	res, err := scripts.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, ErrScriptNotFound
	}

	req.Details.LastModified = lastMod

	return req, nil
}
