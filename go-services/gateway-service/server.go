package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GatewayHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type GatewayServer struct {
	store      Store
	listenAddr string
}

func NewGatewayServer(listenAddr string, store Store) *GatewayServer {
	return &GatewayServer{
		store:      store,
		listenAddr: listenAddr,
	}
}

func (s *GatewayServer) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthz", s.makeHandlerFunc(s.handleHealth))
	router.HandleFunc("POST /login", s.makeHandlerFunc(s.handleLogin))

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

func (s *GatewayServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "OK")
}

func (s *GatewayServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.ContentLength == 0 {
		return fmt.Errorf("request body is empty")
	}
	// Call the auth service to login the user
	resp, err := http.Post("http://localhost:8080/login", "json", r.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login: auth service returned [%s] status code.", resp.Status)
	}
	log.Printf("login successful")
	// Return the response from the auth service
	data := make(map[string]string)
	if json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	return WriteJSON(w, resp.StatusCode, data)
}

func (s *GatewayServer) makeHandlerFunc(f GatewayHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("error: %s", err)
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

// User represents a user in the system
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
