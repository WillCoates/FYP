package business

import "github.com/WillCoates/FYP/common/model"

type SensorReading struct {
	model.SensorData
	Measurement     string
	MeasurementUnit string
	UnitID          string
	UnitName        string
	SensorName      string
}
