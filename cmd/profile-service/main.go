package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Profile struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Language string `json:"language"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	// TODO: Replace with real profile lookup
	profile := Profile{UserID: "u123", Email: "user@example.com", Phone: "+251900000000", Language: "en"}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: profile})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/profile", getProfileHandler)
	fmt.Println("profile-service running on :8087")
	http.ListenAndServe(":8087", nil)
}
