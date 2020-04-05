package service

import (
	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (service *ScriptingService) GetScripts(req *proto.GetScriptsRequest, srv proto.ScriptingService_GetScriptsServer) error {
	scripts := service.db.Collection("scripts")

	_, perms, ok := auth.FromContext(srv.Context())
	if !ok {
		return ErrNoToken
	}

	var users []string

	for _, perm := range perms {
		if perm.Permission == "manageScripting" {
			users = append(users, perm.For)
		}
	}

	if len(users) == 0 {
		return ErrNoPermission
	}

	userIDs, err := util.StringIDToObjectID(users)
	if err != nil {
		return err
	}

	query := bson.M{
		"user": bson.M{"$in": userIDs},
	}

	cur, err := scripts.Find(srv.Context(), query)
	if err != nil {
		return err
	}

	defer cur.Close(srv.Context())
	for cur.Next(srv.Context()) {
		var script model.Script
		err = cur.Decode(&script)
		if err != nil {
			return err
		}
		var protoScript proto.ScriptDetails
		protoScript.Id = script.ID.Hex()
		protoScript.LastModified = script.LastModified
		protoScript.Name = script.Name
		protoScript.Subscriptions = script.Patterns
		err = srv.Send(&protoScript)
		if err != nil {
			return err
		}
	}

	return nil
}
