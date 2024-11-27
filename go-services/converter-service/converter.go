package main

type Converter struct {
	store *MongoStore
	queue *MessageQueue
}

func NewConverter(store *MongoStore, queue *MessageQueue) *Converter {
	return &Converter{
		store: store,
		queue: queue,
	}
}
