package main

import (
	"context"
	"log"
	"time"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func ScriptError(db *mongo.Database, script *model.Script, message string) {
	errors := db.Collection("script_errors")

	var scriptError model.ScriptError
	scriptError.Script = script.ID
	scriptError.Message = message
	scriptError.Timestamp = time.Now().Unix()

	_, err := errors.InsertOne(context.Background(), &scriptError)

	if err != nil {
		log.Println("Failed to record script error", err)
	}
}
