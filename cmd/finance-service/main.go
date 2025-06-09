package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Budget struct {
	UserID string  `json:"user_id"`
	Limit  float64 `json:"limit"`
	Spent  float64 `json:"spent"`
}

type Transaction struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"` // income/expense
	Category string  `json:"category"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

var budgets = []Budget{{UserID: "u123", Limit: 5000, Spent: 1200}}
var transactions = []Transaction{
	{ID: "t1", UserID: "u123", Amount: 100, Type: "expense", Category: "food"},
	{ID: "t2", UserID: "u123", Amount: 200, Type: "income", Category: "salary"},
}

func getBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	for _, b := range budgets {
		if b.UserID == userID {
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: b})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Budget not found"})
}

func listTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userTxs []Transaction
	for _, t := range transactions {
		if t.UserID == userID {
			userTxs = append(userTxs, t)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userTxs})
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	tx.ID = fmt.Sprintf("t%d", len(transactions)+1)
	transactions = append(transactions, tx)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: tx})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/finance/budget", getBudgetHandler)
	http.HandleFunc("/finance/transactions", listTransactionsHandler)
	http.HandleFunc("/finance/transaction", addTransactionHandler)
	fmt.Println("finance-service running on :8089")
	http.ListenAndServe(":8089", nil)
}
