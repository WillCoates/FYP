package service

import (
	"log"

	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	sensorQuery := make(bson.M)
	if len(req.Unit) > 0 {
		sensorQuery["unitid"] = bson.M{"$in": req.Unit}
	}

	if len(req.Sensor) > 0 {
		sensorQuery["info.sensor"] = bson.M{"$in": req.Sensor}
	}

	if len(req.Site) > 0 {
		log.Println("Sites", req.Site)
		var siteIds []primitive.ObjectID
		hasNone := false
		for _, site := range req.Site {
			if site == "None" {
				hasNone = true
			} else {
				siteId, err := service.logic.GetSiteId(users, site, false)
				if err != nil {
					return err
				}
				siteIds = append(siteIds, siteId)
			}
		}
		if !hasNone && len(siteIds) > 0 {
			sensorQuery["site"] = bson.M{"$in": siteIds}
		} else if hasNone && len(siteIds) > 0 {
			sensorQuery["$or"] = bson.A{bson.M{"site": bson.M{"$in": siteIds}}, bson.M{"site": nil}}
		} else {
			sensorQuery["site"] = nil
		}
	}

	readings, err := service.logic.GetSensorReadingsQuery(server.Context(), users, sensorQuery, req.Since)
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
