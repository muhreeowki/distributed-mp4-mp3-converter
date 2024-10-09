package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	store      Store
	listenAddr string
}

func NewServer(listenAddr string, store Store) *Server {
	return &Server{
		store:      store,
		listenAddr: listenAddr,
	}
}

func (s *Server) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("/login", s.handleLogin)

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

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

func CreateJWT(user *User) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
