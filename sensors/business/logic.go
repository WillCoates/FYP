package business

import "go.mongodb.org/mongo-driver/mongo"

type Logic struct {
	db *mongo.Database
}

func NewLogic(db *mongo.Database) *Logic {
	logic := new(Logic)
	logic.db = db
	return logic
}
