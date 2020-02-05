package main

type Config struct {
	AuthService string `toml:"auth-service"`
	KeysURL     string `toml:"keys-url"`
	Binding     string `toml:"binding"`
	MongoURL    string `toml:"mongo-url"`
	MongoDB     string `toml:"mongo-db"`
}
