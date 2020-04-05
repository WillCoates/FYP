package main

import (
	"encoding/json"

	lua "github.com/Shopify/go-lua"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExecScript(db *mongo.Database, ch *amqp.Channel, uid primitive.ObjectID, topic string, payload string) (error, bool) {
	scripts, err := FindScripts(db, uid, topic)
	if err != nil {
		return err, true
	}
	pub := Publish(ch, &uid)

	isJSON := true
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(payload), &data)
	if err != nil {
		isJSON = false
	}

	for script := range scripts {
		l := lua.NewState()
		lua.Require(l, "base", lua.BaseOpen, true)
		lua.Require(l, "bit32", lua.Bit32Open, true)
		lua.Require(l, "math", lua.MathOpen, true)
		lua.Require(l, "string", lua.StringOpen, true)
		lua.Require(l, "table", lua.TableOpen, true)
		l.PushGoFunction(pub)
		l.SetGlobal("publish")
		if isJSON {
			WriteTable(data, l)
		} else {
			l.PushString(payload)
		}
		l.SetGlobal("data")
		err = lua.DoString(l, script.Source)
		if err != nil {
			ScriptError(db, &script, err.Error())
		}
	}
	return nil, false
}
