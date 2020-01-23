package main

import "github.com/WillCoates/FYP/common/model"

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
	MessageId string `json:"msg_id"`
	UnitId    string `json:"unit_id"`
	Sensor    string `json:"sensor"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

func (this *SensorPayload) ToSensorData() (data model.SensorData) {
	data.MessageId = this.MessageId
	data.UnitId = this.UnitId
	data.Sensor = this.Sensor
	data.Value = this.Value
	data.Timestamp = this.Timestamp
	return
}
