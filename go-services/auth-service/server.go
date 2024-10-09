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
	router.HandleFunc("GET /health", s.handleHealth)
	router.HandleFunc("POST /login", s.handleLogin)
	router.HandleFunc("POST /validate", s.handleValidate)

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, "OK")
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

func (s *Server) handleValidate(w http.ResponseWriter, r *http.Request) {
	// Get the token from the request
	token := r.Header.Get("Authorization")
	if token == "" {
		WriteJSON(w, http.StatusUnauthorized, "missing token")
		return
	}

	// Check if the token is a bearer token
	if string(token[:7]) != "Bearer " {
		WriteJSON(w, http.StatusUnauthorized, "invalid authorization header")
		return
	}
	token = token[7:]

	// Verify the token
	_, err := VerifyJWT(token)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, "valid token")
}
