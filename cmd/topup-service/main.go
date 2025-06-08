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

func topupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req TopupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// TODO: Add top-up logic here
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"topup_id": "top123"}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/topup", topupHandler)
	fmt.Println("topup-service running on :8085")
	http.ListenAndServe(":8085", nil)
}
