package main

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
