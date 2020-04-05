package main

// Config holds configuration parameters for the relay service
type Config struct {
	MongoURL     string
	AmqpURL      string
	Listeners    int
	AuthUsername string
	AuthPassword string
}
