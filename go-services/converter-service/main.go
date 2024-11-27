package main

import (
	"log"
)

func main() {
	store, err := NewMongoStore("mongodb://localhost:27017/")
	failOnError(err, "failed to get a MongoStore instance")

	mq, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to get a RabbitMQ instance")

	q, err := mq.channel.QueueDeclare(
		"videoMQ",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare a RabbitMQ queue")

	msgs, err := mq.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "failed to consume a RabbitMQ message: %v")

	var forever chan struct{}

	go func() {
		// Convert the video
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if err := ConvertVideo(store, string(d.Body), "admin"); err != nil {
				log.Printf("failed to convert a video: %v", err)
			}
		}
	}()

	log.Printf(" Converter is Waiting for videos to convert...")
	<-forever
}

func ConvertVideo(store Store, id string, username string) error {
	video, err := store.GetVideoFile(id)
	if err != nil {
		return err
	}
	buf := make([]byte, 2040)
	video.Read(buf)
	log.Printf("Converted a Video: %v", buf)

	return nil
}
