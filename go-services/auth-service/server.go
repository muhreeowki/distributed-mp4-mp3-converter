package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Server represents the HTTP server instance for the auth service
type Server struct {
	store      Store
	listenAddr string
}

// NewServer creates a new Server instance
func NewServer(listenAddr string, store Store) *Server {
	return &Server{
		store:      store,
		listenAddr: listenAddr,
	}
}

// ListenAndServe starts the HTTP server and listens for incoming requests
func (s *Server) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("/login", s.handleLogin)

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

// handleLogin handles the login request and returns a JWT token if the user is valid
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Get the user from the request body
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the user exists in the store
	dbUser, err := s.store.GetUser(user.Email)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Chech if the password is correct
	if dbUser.Password != user.Password {
		log.Printf("User %s failed to log in", user.Email)
		WriteJSON(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Return a success message and a jwt token
	token, err := CreateJWT(user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User %s logged in with token: %s\n", user.Email, token)
	WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
