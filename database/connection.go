package database

import (
	"context"
	"crud-app/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func Connect() {
	uri := config.Config.Mongo.MONGO_URI
	if uri == "" {
		log.Fatal("MONGO_URI not set")
	}

	dbName := config.Config.Mongo.MONGO_DB
	if dbName == "" {
		log.Fatal("MONGO_DB not set")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	MongoClient = client
	MongoDatabase = client.Database(dbName)

	log.Println("MongoDB connected.")
}
