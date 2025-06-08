package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TransferRequest struct {
	SenderID    string  `json:"sender_id"`
	RecipientID string  `json:"recipient_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // CBE, OtherBank, Wallet
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// TODO: Add transfer logic here
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"transaction_id": "tx123"}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/transfer", transferHandler)
	fmt.Println("transfer-service running on :8082")
	http.ListenAndServe(":8082", nil)
}
