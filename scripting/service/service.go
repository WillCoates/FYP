package service

import (
	proto "github.com/WillCoates/FYP/common/protocol/scripting"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScriptingService struct {
	proto.UnimplementedScriptingServiceServer
	db *mongo.Database
}

func NewScriptingService(db *mongo.Database) *ScriptingService {
	service := new(ScriptingService)
	service.db = db
	return service
}
