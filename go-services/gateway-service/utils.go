package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSON writes a JSON response to the http.ResponseWriter with the given status code
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("error encoding JSON response: %w", err)
	}
	return nil
}
