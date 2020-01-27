package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/WillCoates/FYP/auth/business"
	proto "github.com/WillCoates/FYP/common/protocol/auth"
	toml "github.com/pelletier/go-toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func hostGRPC(endpoint string, logic *business.Logic) {
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Println("Failed to start GRPC server")
		log.Println(err)
		return
	}

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, NewAuthService(logic))

	err = server.Serve(lis)
	if err != nil {
		log.Println("Failed running GRPC server")
		log.Println(err)
		return
	}
}

func hostHTTP(endpoint string, logic *business.Logic) {
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Println("Failed to start HTTP server")
		log.Println(err)
		return
	}

	var server http.Server

	server.Handler = CreateRouter()

	err = server.Serve(lis)
	if err != nil {
		log.Println("Failed running HTTP server")
		log.Println(err)
		return
	}
}

func main() {
	config, err := toml.LoadFile("auth.toml")

	if err != nil {
		log.Println("Failed to load config")
		log.Fatalln(err)
	}

	grpcEndpoint := config.GetDefault("grpcEndpoint", ":8081").(string)
	httpEndpoint := config.GetDefault("httpEndpoint", ":8080").(string)
	mongoURL := config.Get("mongoURL").(string)
	currentKey := config.Get("currentKey").(string)
	keysConfig := config.Get("keys").(*toml.Tree)

	keys := make(map[string]*ecdsa.PrivateKey)
	for k, v := range keysConfig.ToMap() {
		data, err := ioutil.ReadFile(v.(string))

		if err != nil {
			log.Println("Failed to load private key ", k)
			log.Fatalln(err)
		}

		pk, err := x509.ParseECPrivateKey(data)

		if err != nil {
			log.Println("Failed to load private key ", k)
			log.Fatalln(err)
		}

		keys[k] = pk
	}

	db, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
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

	logic := business.MakeLogic(db.Database("fyp_auth"))
	logic.SetKeyConfig(currentKey, &keys)

	if grpcEndpoint != "DISABLED" {
		go hostGRPC(grpcEndpoint, logic)
	}

	if httpEndpoint != "DISABLED" {
		go hostHTTP(httpEndpoint, logic)
	}

	forever := make(chan bool)
	<-forever
}
