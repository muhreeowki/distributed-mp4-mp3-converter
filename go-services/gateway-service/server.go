package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GatewayHandlerFunc func(w http.ResponseWriter, r *http.Request) error

// GatewayServer represents the gateway server
type GatewayServer struct {
	store        Store
	messageQueue MessageQueue
	listenAddr   string
}

// NewGatewayServer creates a new GatewayServer
func NewGatewayServer(listenAddr string, store Store, messageQueue MessageQueue) *GatewayServer {
	return &GatewayServer{
		store:        store,
		messageQueue: messageQueue,
		listenAddr:   listenAddr,
	}
}

// ListenAndServe starts the server and listens for incoming requests
func (s *GatewayServer) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthz", s.makeHandlerFunc(s.handleHealth))
	router.HandleFunc("POST /login", s.makeHandlerFunc(s.handleLogin))
	router.HandleFunc("POST /upload", s.makeHandlerFunc(s.handleVideoUpload))

	log.Printf("Server is listening on %s...", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}

// handleHealth handles the health check endpoint
func (s *GatewayServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "OK")
}

// handleLogin handles the login endpoint
func (s *GatewayServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.ContentLength == 0 {
		return fmt.Errorf("request body is empty")
	}
	// Call the auth service to login the user
	resp, err := http.Post(os.Getenv("AUTH_SVC_URL")+"/login", "json", r.Body)
	if err != nil {
		return err
	}

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

// handleVideoUpload handles the video upload endpoint
func (s *GatewayServer) handleVideoUpload(w http.ResponseWriter, r *http.Request) error {
	// Validate the token received in the request header
	if r.Header["Authorization"] == nil {
		return fmt.Errorf("authorization header is missing")
	}
	if err := validateToken(r.Header.Get("Authorization")); err != nil {
		return err
	}
	log.Println("token validated")
	// TODO: Get the users claims from the token. (ie username, email etc.)

	// Parse Video file from request
	if err := r.ParseMultipartForm(20000000); err != nil {
		return fmt.Errorf("failed to parse multipart form: %v", err)
	}

	// retrieve file from form data
	file, handler, err := r.FormFile("mp4File")
	if err != nil {
		return fmt.Errorf("failed to retrive mp4 from request: ", err)
	}
	defer file.Close()

	fmt.Println("File Name:", handler.Filename)
	fmt.Println("File Size:", handler.Size)

	// 1. Store the file in the mongo store using gridfs
	videoId, err := s.store.SaveFile(handler.Filename, file)
	if err != nil {
		return fmt.Errorf("failed to save video file: %v", err)
	}
	log.Printf("Video stored in mongoDB gridfs with ID: %s", videoId)
	// 2. Send a message to the message queue to process the video
	if err := s.messageQueue.SendVideoUploadedMessage(videoId, handler.Size, "bob"); err != nil {
		s.store.DeleteFile(videoId)
		return fmt.Errorf("failed to put video file: %v", err)
	} // TODO: Get users username from the JWT Token

	return WriteJSON(w, http.StatusOK, "upload successful")
}

// TODO: Add download handler

func (s *GatewayServer) makeHandlerFunc(f GatewayHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("error: %v", err)
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

// User represents a user in the system
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// validateToken validates the token by calling the auth service
func validateToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("missing token")
	}

	// Call the auth service to validate the token
	req, err := http.NewRequest("GET", os.Getenv("AUTH_SVC_URL")+"/validate", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token")
	}
	return nil
}
