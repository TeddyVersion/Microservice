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

type Loan struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Term   int     `json:"term_months"`
	Status string  `json:"status"`
}

var loans = []Loan{}

func applyLoanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req Loan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("loan%d", len(loans)+1)
	req.Status = "pending"
	loans = append(loans, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listLoansHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userLoans []Loan
	for _, l := range loans {
		if l.UserID == userID {
			userLoans = append(userLoans, l)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userLoans})
}

func updateLoanStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	for i, l := range loans {
		if l.ID == id {
			loans[i].Status = req.Status
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: loans[i]})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Loan not found"})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/loan/apply", applyLoanHandler)
	http.HandleFunc("/loan/list", listLoansHandler)
	http.HandleFunc("/loan/status", updateLoanStatusHandler)
	fmt.Println("loan-service running on :8088")
	http.ListenAndServe(":8088", nil)
}
