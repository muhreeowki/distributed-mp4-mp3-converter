package main

import "log"

// TODO: Create Server struct
// TODO: Connect to the Mongo DB database and create gridfs bucket
// TODO: Connect to RabbitMQ
//

func main() {
	store := NewMongoStore()

	s := NewGatewayServer(":8080", store)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
