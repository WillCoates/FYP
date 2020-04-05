package service

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service *ScriptingService) DeleteScript(ctx context.Context, req *proto.Script) (*proto.Script, error) {
	scripts := service.db.Collection("scripts")
	errors := service.db.Collection("script_errors")

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

	filter := bson.M{
		"_id":  scriptID,
		"user": bson.M{"$in": userIDs},
	}

	res, err := scripts.DeleteOne(ctx, filter)

	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, ErrScriptNotFound
	}

	filter = bson.M{
		"script": scriptID,
	}

	_, err = errors.DeleteMany(ctx, filter)

	if err != nil {
		return nil, err
	}

	return req, nil
}
