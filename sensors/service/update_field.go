package service

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
)

func (service *SensorsService) UpdateField(ctx context.Context, req *proto.Field) (*proto.Field, error) {
	var err error

	_, perms, ok := auth.FromContext(ctx)
	if !ok {
		return nil, ErrNoToken
	}

	var users []string

	for _, perm := range perms {
		if perm.Permission == "configureSensor" {
			users = append(users, perm.For)
		}
	}

	if len(users) == 0 {
		return nil, ErrNoPermission
	}

	var field model.SensorSite
	field.Crop = req.Crop
	field.Name = req.Name

	err = service.logic.UpdateField(ctx, &field, users)
	if err != nil {
		return nil, err
	}

	return req, nil
}
