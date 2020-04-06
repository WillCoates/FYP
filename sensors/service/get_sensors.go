package service

import (
	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	query := make(bson.M)

	if len(req.Site) > 0 {
		siteIds := make([]primitive.ObjectID, len(req.Site))
		for i, site := range req.Site {
			var err error
			siteIds[i], err = service.logic.GetSiteId(users, site, false)
			if err != nil {
				return err
			}
		}
		query["site"] = bson.M{"$in": siteIds}
	}

	if req.Name != "" {
		query["name"] = req.Name
	}

	if len(req.Unit) > 0 {
		query["unitid"] = bson.M{"$in": req.Unit}
	}

	if len(req.Sensor) > 0 {
		query["info.sensor"] = bson.M{"$in": req.Sensor}
	}

	if !req.IncludeHidden {
		query["info.hidden"] = false
	}

	sensors, err := service.logic.GetSensors(server.Context(), query, users...)
	if err != nil {
		return err
	}

	for sensor := range sensors {
		var info proto.SensorInfo
		info.Name = sensor.Name
		info.Unit = sensor.UnitID
		info.Latitude = sensor.Latitude
		info.Longitude = sensor.Longitude
		info.Sensor = sensor.Info.Sensor
		info.Measurementname = sensor.Info.Measurement
		info.Measurementunit = sensor.Info.Units
		info.Hidden = sensor.Info.Hidden
		info.Site, _ = service.logic.GetSiteName(users, &sensor.Site)
		server.Send(&info)
	}

	return nil
}
