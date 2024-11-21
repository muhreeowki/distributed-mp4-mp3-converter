package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GatewayHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type GatewayServer struct {
	store        Store
	messageQueue MessageQueue
	listenAddr   string
}

func NewGatewayServer(listenAddr string, store Store, messageQueue MessageQueue) *GatewayServer {
	return &GatewayServer{
		store:        store,
		messageQueue: messageQueue,
		listenAddr:   listenAddr,
	}
}

func (s *GatewayServer) ListenAndServe() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /healthz", s.makeHandlerFunc(s.handleHealth))
	router.HandleFunc("POST /login", s.makeHandlerFunc(s.handleLogin))
	router.HandleFunc("POST /upload", s.makeHandlerFunc(s.handleVideoUpload))

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

func (s *GatewayServer) handleVideoUpload(w http.ResponseWriter, r *http.Request) error {
	// Validate the token received in the request header
	if r.Header["Authorization"] == nil {
		return fmt.Errorf("authorization header is missing")
	}
	if err := validateToken(r.Header.Get("Authorization")); err != nil {
		return err
	}
	log.Println("token validated")

	// Parse Video file from request
	r.ParseMultipartForm(10 << 20)

	// retrieve file from form data
	file, handler, err := r.FormFile("mp4File")
	if err != nil {
		return fmt.Errorf("failed to retrive mp4 from request: ", err)
	}
	defer file.Close()

	fmt.Println("File Name:", handler.Filename)
	fmt.Println("File Size:", handler.Size)

	// TODO: Upload the file to the store

	// 1. Store the file in the mongo store using gridfs
	//    vid_id := s.store.SaveFile(file)
	// 2. Send a message to the message queue to process the video
	//    s.messageQueue.SendMessage(vid_id)
	// 3. Return a response

	return WriteJSON(w, http.StatusOK, "upload successful")
}

// TODO: Add download handler

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
