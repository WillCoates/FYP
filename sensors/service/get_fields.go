package service

import (
	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
	"go.mongodb.org/mongo-driver/bson"
)

func (service *SensorsService) GetFields(req *proto.GetFieldsRequest, server proto.SensorsService_GetFieldsServer) error {
	_, perms, ok := auth.FromContext(server.Context())
	if !ok {
		return ErrNoToken
	}

	var users []string

	for _, perm := range perms {
		if perm.Permission == "readSensor" {
			users = append(users, perm.For)
		}
	}

	if len(users) == 0 {
		return ErrNoPermission
	}

	query := make(bson.M)

	if len(req.Name) > 0 {
		query["name"] = bson.M{"$in": req.Name}
	}

	fields, err := service.logic.GetFields(server.Context(), users, query)
	if err != nil {
		return err
	}

	for field := range fields {
		var info proto.Field
		info.Name = field.Name
		info.Crop = field.Crop
		service.logic.GetSensors(server.Context(), bson.M{"site": field.ID}, users...)
		server.Send(&info)
	}

	return nil
}
