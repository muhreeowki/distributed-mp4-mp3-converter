package main

type Store interface{}

type MongoStore struct{}

func NewMongoStore() *MongoStore {
	return &MongoStore{}
}
