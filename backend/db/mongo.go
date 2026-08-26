package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is the shared MongoDB client.
var Client *mongo.Client

// DB is the working database handle.
var DB *mongo.Database

// Connect establishes a connection to MongoDB and pings it to verify.
func Connect(uri, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping error: %v", err)
	}

	Client = client
	DB = client.Database(dbName)
	log.Printf("connected to MongoDB database %q", dbName)
}

// Collection returns a handle to a named collection.
func Collection(name string) *mongo.Collection {
	return DB.Collection(name)
}
