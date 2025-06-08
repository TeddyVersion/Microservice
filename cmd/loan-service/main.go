package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoanApplication struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Term   int     `json:"term_months"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func applyLoanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req LoanApplication
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// TODO: Add loan application logic here
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"application_id": "loan123"}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/loan/apply", applyLoanHandler)
	fmt.Println("loan-service running on :8088")
	http.ListenAndServe(":8088", nil)
}
