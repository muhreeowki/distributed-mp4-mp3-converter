package main

import (
	"log"
)

func main() {
	store, err := NewPostgersStore()
	if err != nil {
		log.Fatal(err)
		return
	}

	server := NewServer(":8080", store)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
