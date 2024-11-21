package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface{}

type MongoStore struct {
	gridfs *gridfs.Bucket
	client *mongo.Client
}

func NewMongoStore(conStr string) (*MongoStore, error) {
	// Set the server API version for the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conStr).SetServerAPIOptions(serverAPI)

	// Create a new client
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Get the database
	db := client.Database("admin")

	// Check the connection
	var result bson.M
	if err := db.RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB.")

	// Create a new GridFS bucket
	gfs, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GridFS bucket: %v", err)
	}

	return &MongoStore{
		gridfs: gfs,
		client: client,
	}, nil
}
