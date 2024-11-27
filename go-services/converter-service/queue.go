package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	SendVideoUploadedMessage(id string, size int64, username string) error
}

type RabbitMQ struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewRabbitMQ(connStr string) (*RabbitMQ, error) {
	// Connect to RabbitMQ Serevr
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: err", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a RabbitMQ channel: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ.")

	return &RabbitMQ{
		channel: ch,
		conn:    conn,
	}, nil
}

func (mq *RabbitMQ) Listen() {
}

func (mq *RabbitMQ) Close() {
	mq.channel.Close()
	mq.conn.Close()
}

func (mq *RabbitMQ) SendVideoUploadedMessage(mp3Id string, videoId string, size int64, username string) error {
	msg := map[string]string{
		"videoId":  videoId,
		"mp3Id":    mp3Id,
		"username": username,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to send a RabbitMQ message: %v", err)
	}

	queue, err := mq.channel.QueueDeclare(
		"mp3Q",
		false,
		false,
		false,
		false,
		nil,
	)

	return mq.channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
}
