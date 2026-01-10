package config

import (
	"context"
	"log"
	"time"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)


func ConnectMongo() {

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL")).SetServerAPIOptions(serverAPI)
	   
	client, err := mongo.Connect(
		opts,
	)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	log.Println("MongoDB connected..")

	Client = client
	DB = client.Database("code_editor")
}

