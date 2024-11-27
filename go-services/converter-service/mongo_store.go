package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetVideoFile(objectId string) (io.ReadCloser, error)
	SaveMP3File(filename string, file io.Reader) (string, error)
	DeleteMP3File(objectId string) error
}

type MongoStore struct {
	gfsVideo *gridfs.Bucket
	gfsMp3   *gridfs.Bucket
	client   *mongo.Client
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

	// Get the databases
	videos_db := client.Database("videos")
	mp3_db := client.Database("mp3")

	// Check the connection for the databases
	var result bson.M
	if err := videos_db.RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Println(result)
		return nil, fmt.Errorf("failed to ping videos database: %v", err)
	}
	if err := mp3_db.RunCommand(context.Background(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Println(result)
		return nil, fmt.Errorf("failed to ping videos database: %v", err)
	}
	log.Println("Successfully connected to Videos and MP3 DBs.")

	// Create a new GridFS bucket
	gfsVideo, err := gridfs.NewBucket(videos_db)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GridFS bucket: %v", err)
	}
	gfsMp3, err := gridfs.NewBucket(videos_db)
	if err != nil {
		return nil, fmt.Errorf("Failed to create GridFS bucket: %v", err)
	}

	return &MongoStore{
		gfsVideo: gfsVideo,
		gfsMp3:   gfsMp3,
		client:   client,
	}, nil
}

func (s *MongoStore) GetVideoFile(objectId string) (io.ReadCloser, error) {
	return s.gfsVideo.OpenDownloadStream(objectId)
}

func (s *MongoStore) SaveMP3File(filename string, file io.Reader) (string, error) {
	objectId, err := s.gfsMp3.UploadFromStream(filename, file)
	if err != nil {
		return "", err
	}
	return objectId.String(), nil
}

func (s *MongoStore) DeleteMP3File(objectId string) error {
	return s.gfsMp3.Delete(objectId)
}
