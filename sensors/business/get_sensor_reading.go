package business

import (
	"context"
	"log"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (logic *Logic) GetSensorReadingsQuery(ctx context.Context, users []string, sensorQuery bson.M, since int64) (chan SensorReading, error) {
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

	result, err := readings.Find(ctx, readingQuery)
	if err != nil {
		return nil, err
	}

	readingsChan := make(chan SensorReading)

	go func() {
		defer result.Close(ctx)
		for result.Next(ctx) {
			var modelReading model.SensorData
			err := result.Decode(&modelReading)
			if err != nil {
				log.Println("Failed to decode sensor reading")
				log.Println(err)
			} else {
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
		close(readingsChan)
	}()

	return readingsChan, nil
}

// GetSensorReadings Deprecated: Use GetSensorReadingsQuery with a BSON query to fetch readings
func (logic *Logic) GetSensorReadings(ctx context.Context, id []string, sensor []string, since int64, users []string) (chan SensorReading, error) {
	sensorQuery := make(bson.M)

	if id != nil {
		sensorQuery["unitid"] = bson.M{"$in": id}
	}

	if sensor != nil {
		sensorQuery["info.sensor"] = bson.M{"$in": sensor}
	}

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

	result, err := readings.Find(ctx, readingQuery)
	if err != nil {
		return nil, err
	}

	readingsChan := make(chan SensorReading)

	go func() {
		defer result.Close(ctx)
		for result.Next(ctx) {
			var modelReading model.SensorData
			err := result.Decode(&modelReading)
			if err != nil {
				log.Println("Failed to decode sensor reading")
				log.Println(err)
			} else {
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
		close(readingsChan)
	}()

	return readingsChan, nil
}
