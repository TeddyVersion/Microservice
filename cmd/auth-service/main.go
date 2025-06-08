package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// In-memory store for demonstration
var users = []RegisterRequest{}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// Simulate storing the user
	users = append(users, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"user_id": fmt.Sprintf("u%d", len(users))}})
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: users})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/users", listUsersHandler)
	fmt.Println("auth-service running on :8080")
	http.ListenAndServe(":8080", nil)
}
