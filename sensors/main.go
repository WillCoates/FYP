package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/WillCoates/FYP/common/auth"
	proto "github.com/WillCoates/FYP/common/protocol/sensors"
	"github.com/WillCoates/FYP/sensors/business"
	"github.com/WillCoates/FYP/sensors/service"
	"github.com/pelletier/go-toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	configFile := "sensors.toml"
	var config Config

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	configRaw, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Println("Failed to read configuration file")
		log.Fatalln(err)
	}

	err = toml.Unmarshal(configRaw, &config)
	if err != nil {
		log.Println("Failed to unmarshal configuration file")
		log.Fatalln(err)
	}

	authClient, err := auth.NewAuthClient(config.AuthService, config.KeysURL)

	if err != nil {
		log.Println("Failed to authentication service client")
		log.Fatalln(err)
	}

	db, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		log.Println("Failed to create MongoDB client")
		log.Fatalln(err)
	}

	dbctx, dbcancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer dbcancel()

	err = db.Connect(dbctx)
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		log.Fatalln(err)
	}

	logic := business.NewLogic(db.Database(config.MongoDB))

	server := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryServerInteceptor(nil, authClient)),
		grpc.StreamInterceptor(auth.StreamServerInteceptor(nil, authClient)))

	proto.RegisterSensorsServiceServer(server, service.NewSensorsService(logic))

	lis, err := net.Listen("tcp", config.Binding)
	if err != nil {
		log.Println("Failed to bind socket")
		log.Fatalln(err)
	}

	err = server.Serve(lis)
	if err != nil {
		log.Println("Failed to serve service")
		log.Fatalln(err)
	}
}
