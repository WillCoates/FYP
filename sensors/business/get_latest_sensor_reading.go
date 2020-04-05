package business

import (
	"context"
	"log"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (logic *Logic) GetLatestSensorReadingsQuery(ctx context.Context, users []string, sensorQuery bson.M, since int64) (chan SensorReading, error) {
	sensorChannel, err := logic.GetSensors(ctx, sensorQuery, users...)

	if err != nil {
		return nil, err
	}

	var sensorIds []primitive.ObjectID
	var sensorMapping = make(map[primitive.ObjectID]model.Sensor)
	for sen := range sensorChannel {
		sensorIds = append(sensorIds, sen.ID)
		sensorMapping[sen.ID] = sen
	}

	var readings = logic.db.Collection("sensor_readings")
	readingQuery := make(bson.M)

	readingQuery["sensor"] = bson.M{"$in": sensorIds}
	readingQuery["timestamp"] = bson.M{"$gte": since}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"timestamp": -1})

	result, err := readings.Find(ctx, readingQuery, findOptions)
	if err != nil {
		return nil, err
	}

	readingsChan := make(chan SensorReading)

	go func() {
		found := make(map[primitive.ObjectID]bool)
		defer result.Close(ctx)
		for result.Next(ctx) {
			var modelReading model.SensorData
			err := result.Decode(&modelReading)
			if err != nil {
				log.Println("Failed to decode sensor reading")
				log.Println(err)
			} else {
				_, exists := found[modelReading.Sensor]
				if !exists {
					found[modelReading.Sensor] = true
					var reading SensorReading
					reading.SensorData = modelReading
					reading.Measurement = sensorMapping[modelReading.Sensor].Info.Measurement
					reading.MeasurementUnit = sensorMapping[modelReading.Sensor].Info.Units
					reading.UnitID = sensorMapping[modelReading.Sensor].UnitID
					reading.UnitName = sensorMapping[modelReading.Sensor].Name
					reading.SensorName = sensorMapping[modelReading.Sensor].Info.Sensor
					readingsChan <- reading
				}
			}
		}
		close(readingsChan)
	}()

	return readingsChan, nil
}
