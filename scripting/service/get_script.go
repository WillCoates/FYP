package service

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service *ScriptingService) GetScript(ctx context.Context, req *proto.GetScriptRequest) (*proto.Script, error) {
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

	scriptID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"user": bson.M{"$in": userIDs},
		"_id":  scriptID,
	}

	var script model.Script

	err = scripts.FindOne(ctx, query).Decode(&script)
	if err != nil {
		return nil, err
	}

	protoScript := new(proto.Script)
	protoScript.Source = script.Source
	protoScript.Details = new(proto.ScriptDetails)
	protoScript.Details.Id = script.ID.Hex()
	protoScript.Details.LastModified = script.LastModified
	protoScript.Details.Name = script.Name
	protoScript.Details.Subscriptions = script.Patterns

	return protoScript, nil
}
