package business

import (
	"context"

	"github.com/WillCoates/FYP/common/model"
	"github.com/WillCoates/FYP/common/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (logic *Logic) DeleteSensor(ctx context.Context, sensor *model.Sensor, users []string) error {
	sensors := logic.db.Collection("sensors")
	readings := logic.db.Collection("sensor_readings")

	userIds, err := util.StringIDToObjectID(users)

	query := make(bson.M)

	query["user"] = bson.M{"$in": userIds}
	query["unitid"] = sensor.UnitID
	query["info.sensor"] = sensor.Info.Sensor

	var deletedSensor model.Sensor
	err = sensors.FindOneAndDelete(ctx, query).Decode(&deletedSensor)
	if err != nil {
		return err
	}

	readings.DeleteMany(ctx, bson.M{"sensor": deletedSensor.ID})

	return nil
}
