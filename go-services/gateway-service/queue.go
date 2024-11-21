package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface{}

type RabbitMQ struct {
	channel *amqp.Channel
}

func NewRabbitMQ(connStr string) (*RabbitMQ, error) {
	// Connect to RabbitMQ Serevr
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: err", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	log.Println("Successfully connected to RabbitMQ.")

	return &RabbitMQ{
		channel: ch,
	}, nil
}
