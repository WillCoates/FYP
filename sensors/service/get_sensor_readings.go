package service

import (
	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
)

func (service *SensorsService) GetSensorReadings(req *proto.GetSensorReadingsRequest, server proto.SensorsService_GetSensorReadingsServer) error {
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

	readings, err := service.logic.GetSensorReadings(server.Context(), req.Unit, req.Sensor, req.Since, users)
	if err != nil {
		return err
	}

	for reading := range readings {
		var info proto.SensorData
		info.Reading = reading.Value
		info.Measurementname = reading.Measurement
		info.Measurementunit = reading.MeasurementUnit
		info.Sensor = reading.SensorName
		info.Timestamp = reading.Timestamp
		info.Unit = reading.UnitID
		info.UnitName = reading.UnitName
		server.Send(&info)
	}

	return nil
}
