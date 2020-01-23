package main

import (
	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Meshlium payload:
{
	"msg_id": "#ID#",
	"unit_id": "#ID_WASP#",
	"sensor": "#SENSOR#",
	"value": "#VALUE#",
	"timestamp": #TS('U')
}
*/

type SensorPayload struct {
	MessageID string `json:"msg_id"`
	UnitID    string `json:"unit_id"`
	Sensor    string `json:"sensor"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func (this *SensorPayload) ToSensorData() (data model.SensorData) {
	data.ID = primitive.NewObjectID()
	data.MessageID = this.MessageID
	data.UnitID = this.UnitID
	data.Sensor = this.Sensor
	data.Value = this.Value
	data.Timestamp = this.Timestamp
	return
}
