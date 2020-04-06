package service

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
)

func (service *SensorsService) DeleteSensor(ctx context.Context, req *proto.SensorInfo) (*proto.SensorInfo, error) {
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

	var info model.SensorInfo
	info.Sensor = req.Sensor

	var sensor model.Sensor
	sensor.UnitID = req.Unit
	sensor.Info = &info

	err = service.logic.DeleteSensor(ctx, &sensor, users)
	if err != nil {
		return nil, err
	}

	return req, nil
}
