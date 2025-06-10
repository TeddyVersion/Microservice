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

var profiles = []Profile{
	{UserID: "u123", Email: "user@example.com", Phone: "+251900000000", Language: "en"},
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
	userID := r.URL.Query().Get("user_id")
	for _, p := range profiles {
		if p.UserID == userID {
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: p})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Profile not found"})
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	for i, p := range profiles {
		if p.UserID == req.UserID {
			profiles[i] = req
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
			return
		}
	}
	profiles = append(profiles, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func deleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	for i, p := range profiles {
		if p.UserID == userID {
			profiles = append(profiles[:i], profiles[i+1:]...)
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Message: "Profile deleted"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Profile not found"})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/profile", getProfileHandler)
	http.HandleFunc("/profile/update", updateProfileHandler)
	http.HandleFunc("/profile/delete", deleteProfileHandler)
	fmt.Println("profile-service running on :8087")
	http.ListenAndServe(":8087", nil)
}
