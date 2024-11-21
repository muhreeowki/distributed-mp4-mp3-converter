package main

// TODO: Create Server struct
// TODO: Connect to the Mongo DB database and create gridfs bucket
// TODO: Connect to RabbitMQ
//

func main() {
	store, err := NewMongoStore("mongodb://localhost:27017/")
	failOnError(err, "")

	queue, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	failOnError(err, "")

	server := NewGatewayServer(":3000", store, queue)
	err = server.ListenAndServe()
	failOnError(err, "Failed to start server")
}
