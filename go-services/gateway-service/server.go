package main

import (
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

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

func (s *GatewayServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "OK")
}

func (s *GatewayServer) makeHandlerFunc(f GatewayHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("error: %s", err)
		}
	}
}
