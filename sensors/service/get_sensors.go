package service

import (
	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
)

func (service *SensorsService) GetSensors(req *proto.GetSensorsRequest, server proto.SensorsService_GetSensorsServer) error {
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

	sensors, err := service.logic.GetSensors(server.Context(), nil, users...)
	if err != nil {
		return err
	}

	for sensor := range sensors {
		var info proto.SensorInfo
		info.Name = sensor.Name
		info.Unit = sensor.UnitID
		info.Sensor = sensor.Info.Sensor
		info.Measurementname = sensor.Info.Measurement
		info.Measurementunit = sensor.Info.Units
		info.Hidden = sensor.Info.Hidden
		server.Send(&info)
	}

	return nil
}
