package service

import (
	"context"
	"time"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service *ScriptingService) AddScript(ctx context.Context, req *proto.Script) (*proto.Script, error) {
	scripts := service.db.Collection("scripts")

	tkn, perms, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	hasPermission := false

	for _, perm := range perms {
		if perm.Permission == "manageScripting" {
			if perm.For == tkn.Payload.Subject {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return nil, ErrNoPermission
	}

	user, err := primitive.ObjectIDFromHex(tkn.Payload.Subject)

	if err != nil {
		return nil, err
	}

	lastMod := time.Now().Unix()

	var script model.Script
	script.User = user
	script.LastModified = lastMod
	script.Source = req.Source
	script.Name = req.Details.Name
	script.Patterns = req.Details.Subscriptions

	res, err := scripts.InsertOne(ctx, &script)
	if err != nil {
		return nil, err
	}

	req.Details.Id = res.InsertedID.(primitive.ObjectID).Hex()
	req.Details.LastModified = lastMod

	return req, nil
}
