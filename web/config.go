package main

type Config struct {
	RedisURL    string
	TemplateDir string
	Global      map[string]string
}
