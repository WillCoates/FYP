package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SensorData struct {
	Id primitive.ObjectID `bson:'_id'`
	User string `bson:'user'`
	MessageId string `bson:'messageid'`
	UnitId    string `bson:'unitid'`
	Sensor    string `bson:'sensor'`
	Value     string `bson:'value'`
	Timestamp int64 `bson:'timestamp'`
}
