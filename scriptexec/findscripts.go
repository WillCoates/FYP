package main

import (
	"context"
	"log"

	"github.com/WillCoates/FYP/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindScripts(db *mongo.Database, uid primitive.ObjectID, topic string) (chan model.Script, error) {
	scripts := db.Collection("scripts")
	ctx := context.Background()
	// TODO: Optimize this more
	cursor, err := scripts.Find(ctx, bson.M{"user": uid})
	if err != nil {
		return nil, err
	}

	scriptChan := make(chan model.Script)
	go func() {
		defer close(scriptChan)
		for cursor.Next(ctx) {
			var script model.Script
			err := cursor.Decode(&script)
			if err != nil {
				log.Println("Failed to decode script", err)
				continue
			}

			for _, pattern := range script.Patterns {
				regex, err := TopicToRegex(pattern)
				if err != nil {
					ScriptError(db, &script, "Badly formatted topic: "+pattern)
					log.Println("Failed to convert topic to regex", err)
					continue
				}
				if regex.MatchString(topic) {
					scriptChan <- script
					break
				}
			}
		}
	}()

	return scriptChan, nil
}
