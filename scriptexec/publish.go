package main

import (
	"encoding/json"
	"log"

	lua "github.com/Shopify/go-lua"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Publish(ch *amqp.Channel, userID *primitive.ObjectID) func(state *lua.State) int {
	topicPrefix := userID.Hex() + "."
	return func(l *lua.State) int {
		var err error
		topic, ok := l.ToString(1)
		if !ok {
			l.PushString("bad argument #1 to 'publish' (string expected, got " + l.TypeOf(1).String() + ")")
			l.Error()
			return 0
		}
		var msg amqp.Publishing
		switch {
		case l.IsString(2):
			payload, _ := l.ToString(2)
			msg.Body = []byte(payload)
			msg.ContentType = "text/plain"
		case l.IsTable(2):
			payload, ok := ReadTable(2, l)
			if !ok {
				l.PushString("argument #2 to 'publish' failed to read table")
				l.Error()
				return 0
			}
			msg.Body, err = json.Marshal(&payload)
			if err != nil {
				log.Println("Publish failed to marshal", err)
				l.PushString("argument #2 to 'publish' failed to marshal to JSON")
				l.Error()
				return 0
			}
			msg.ContentType = "application/json"
		default:
			l.PushString("bad argument #2 to 'publish' (string or table expected, got " + l.TypeOf(1).String() + ")")
			l.Error()
			return 0
		}

		msg.AppId = "scriptexec"

		ch.Publish("amq.topic", topicPrefix+topic, false, false, msg)
		return 1
	}
}
