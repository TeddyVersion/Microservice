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

type Transfer struct {
	ID          string  `json:"id"`
	SenderID    string  `json:"sender_id"`
	RecipientID string  `json:"recipient_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // CBE, OtherBank, Wallet
	Status      string  `json:"status"`
}

var transfers = []Transfer{}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req Transfer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("tx%d", len(transfers)+1)
	req.Status = "completed"
	transfers = append(transfers, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listTransfersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userTxs []Transfer
	for _, t := range transfers {
		if t.SenderID == userID || t.RecipientID == userID {
			userTxs = append(userTxs, t)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userTxs})
}

func getTransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	for _, t := range transfers {
		if t.ID == id {
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: t})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Transfer not found"})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/transfer", transferHandler)
	http.HandleFunc("/transfers", listTransfersHandler)
	http.HandleFunc("/transfer/get", getTransferHandler)
	fmt.Println("transfer-service running on :8082")
	http.ListenAndServe(":8082", nil)
}
