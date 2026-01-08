package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)


func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(""),
	)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
	}

	log.Println("✅ MongoDB connected")

	// Assign globals
	Client = client
	DB = client.Database("code_editor")
}
