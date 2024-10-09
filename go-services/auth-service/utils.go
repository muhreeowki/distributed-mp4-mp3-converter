package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthClaims represents the JWT claims for the auth service
type AuthClaims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool
	jwt.RegisteredClaims
}

// CreateJWT creates a new JWT token for the user
func CreateJWT(user *User) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Create the Claims
	claims := &AuthClaims{
		Email:    user.Email,
		Password: user.Password,
		Admin:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 3)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT verifies the JWT token and returns the token if it is valid
func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	// Check if there was an error parsing the token
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*AuthClaims); ok {
		log.Println(claims.Email, claims.IssuedAt.Time)
		return token, nil
	} else {
		return nil, fmt.Errorf("unknown claims type: %+v", claims)
	}
}

// WriteJSON writes a JSON response to the http.ResponseWriter with the given status code
func WriteJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
