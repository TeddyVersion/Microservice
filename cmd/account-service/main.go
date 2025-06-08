package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
	Type    string  `json:"type"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// In-memory store for demonstration
var accounts = []Account{
	{ID: "acc1", Balance: 1000.0, Type: "savings"},
	{ID: "acc2", Balance: 250.5, Type: "checking"},
}

func getAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: accounts})
}

func getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	for _, acc := range accounts {
		if acc.ID == id {
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: acc})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Account not found"})
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var acc Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	acc.ID = fmt.Sprintf("acc%d", len(accounts)+1)
	accounts = append(accounts, acc)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: acc})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAccountsHandler(w, r)
		case http.MethodPost:
			createAccountHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		}
	})
	http.HandleFunc("/account", getAccountByIDHandler) // GET /account?id=acc1
	fmt.Println("account-service running on :8081")
	http.ListenAndServe(":8081", nil)
}
