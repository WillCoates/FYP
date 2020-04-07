package service

import (
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (service *ScriptingService) GetScriptErrors(req *proto.GetScriptErrorsRequest, srv proto.ScriptingService_GetScriptErrorsServer) error {
	scripts := service.db.Collection("scripts")
	errors := service.db.Collection("script_errors")

	query := make(bson.M)
	if len(req.Id) > 0 {
		ids, err := util.StringIDToObjectID(req.Id)
		if err != nil {
			return err
		}
		query["script"] = bson.M{"$in": ids}
	}
	query["timestamp"] = bson.M{"$gte": req.Since}

	opts := options.Find()
	opts.SetSort(bson.M{"timestamp": -1})
	if req.Limit > 0 {
		opts.SetLimit(int64(req.Limit))
	}

	cur, err := errors.Find(srv.Context(), query, opts)
	if err != nil {
		return err
	}

	for cur.Next(srv.Context()) {
		var scriptErr model.ScriptError
		var script model.Script

		err = cur.Decode(&scriptErr)
		if err != nil {
			return err
		}

		err = scripts.FindOne(srv.Context(), bson.M{"_id": scriptErr.Script}).Decode(&script)

		var protoErr proto.ScriptError
		protoErr.Message = scriptErr.Message
		protoErr.Timestamp = scriptErr.Timestamp
		protoErr.Script = new(proto.ScriptDetails)
		protoErr.Script.Id = script.ID.Hex()
		protoErr.Script.LastModified = script.LastModified
		protoErr.Script.Name = script.Name
		protoErr.Script.Subscriptions = script.Patterns

		err = srv.Send(&protoErr)
		if err != nil {
			return err
		}
	}

	return nil
}
