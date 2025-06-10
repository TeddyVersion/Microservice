package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TopupRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"` // airtime or data
	Phone  string  `json:"phone"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type Topup struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"` // airtime or data
	Phone  string  `json:"phone"`
	Status string  `json:"status"`
}

var topups = []Topup{}

func topupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req Topup
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("top%d", len(topups)+1)
	req.Status = "completed"
	topups = append(topups, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listTopupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userTopups []Topup
	for _, t := range topups {
		if t.UserID == userID {
			userTopups = append(userTopups, t)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userTopups})
}

func getTopupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	for _, t := range topups {
		if t.ID == id {
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: t})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Topup not found"})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/topup", topupHandler)
	http.HandleFunc("/topups", listTopupsHandler)
	http.HandleFunc("/topup/get", getTopupHandler)
	fmt.Println("topup-service running on :8085")
	http.ListenAndServe(":8085", nil)
}
