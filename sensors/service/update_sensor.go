package service

import (
	"context"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/common/model"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
)

func (service *SensorsService) UpdateSensor(ctx context.Context, req *proto.SensorInfo) (*proto.SensorInfo, error) {
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
	info.Hidden = req.GetHidden()
	info.Measurement = req.GetMeasurementname()
	info.Units = req.GetMeasurementunit()
	info.Sensor = req.GetSensor()

	var sensor model.Sensor
	sensor.Name = req.GetName()
	sensor.UnitID = req.GetUnit()
	sensor.Latitude = req.GetLatitude()
	sensor.Longitude = req.GetLongitude()
	sensor.Info = &info
	sensor.Site, err = service.logic.GetSiteId(users, req.Site, true)
	if err != nil {
		return nil, err
	}

	err = service.logic.UpdateSensor(ctx, &sensor, users)
	if err != nil {
		return nil, err
	}

	return req, nil
}
