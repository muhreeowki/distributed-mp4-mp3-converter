package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// AuthServer represents the HTTP server instance for the auth service
type AuthServer struct {
	store      Store
	listenAddr string
}

// NewAuthServer creates a new Server instance
func NewAuthServer(listenAddr string, store Store) *AuthServer {
	return &AuthServer{
		store:      store,
		listenAddr: listenAddr,
	}
}

// ListenAndServe starts the HTTP server and listens for incoming requests
func (s *AuthServer) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthz", s.handleHealth)
	router.HandleFunc("POST /login", s.handleLogin)
	router.HandleFunc("GET /validate", s.handleValidate)

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

func (s *AuthServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, "OK")
}

// handleLogin handles the login request and returns a JWT token if the user is valid
func (s *AuthServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Get the user from the request body
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the user has provided the email and password
	if user.Email == "" || user.Password == "" {
		WriteJSON(w, http.StatusInternalServerError, fmt.Errorf("missing credentials"))
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

func (s *AuthServer) handleValidate(w http.ResponseWriter, r *http.Request) {
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
