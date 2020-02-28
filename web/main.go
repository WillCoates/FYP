package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/WillCoates/FYP/common/auth"
	"github.com/WillCoates/FYP/web/business"
	"github.com/WillCoates/FYP/web/framework"
	"github.com/pelletier/go-toml"
)

func main() {
	configFile := "web.toml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	var config Config

	configTree, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	err = configTree.Unmarshal(&config)
	if err != nil {
		log.Fatalln("Failed to unmarshal config", err)
	}

	sessionManager, err := framework.NewSessionManager(config.RedisURL)
	if err != nil {
		log.Fatalln("Failed to create session manager", err)
	}

	templateManager, err := framework.NewTemplateManager(config.TemplateDir)
	if err != nil {
		log.Fatalln("Failed to create template manager", err)
	}

	logic := new(business.Logic)
	logic.MasterKey = auth.MasterKey()
	logic.Config = config.Global

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to start HTTP server", err)
	}

	var server http.Server

	server.Handler = CreateRouter(logic, sessionManager, templateManager)

	err = server.Serve(lis)
	if err != nil {
		log.Println("Failed running HTTP server")
		log.Println(err)
		return
	}
}
