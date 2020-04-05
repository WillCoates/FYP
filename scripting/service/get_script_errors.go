package service

import (
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *ScriptingService) GetScriptErrors(req *proto.GetScriptErrorsRequest, srv proto.ScriptingService_GetScriptErrorsServer) error {
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
	opts.SetSort(bson.M{"timestamp": 1})

	errors.Find(srv.Context(), query, opts)

	return status.Errorf(codes.Unimplemented, "method GetScriptErrors not implemented")
}
