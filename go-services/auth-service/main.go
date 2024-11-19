package main

import (
	"log"
)

func main() {
	// Create a new PostgresStore instance
	store, err := NewPostgersStore()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new Server instance
	server := NewAuthServer(":8080", store)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
