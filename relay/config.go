package main

// Config holds configuration parameters for the relay service
type Config struct {
	MongoURL     string
	AmqpURL      string
	Listeners    int
	AuthService  string
	AuthUsername string
	AuthPassword string
}
