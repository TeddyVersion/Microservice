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

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func getBudgetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	// TODO: Replace with real budget lookup
	budget := Budget{UserID: "u123", Limit: 5000, Spent: 1200}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: budget})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/finance/budget", getBudgetHandler)
	fmt.Println("finance-service running on :8089")
	http.ListenAndServe(":8089", nil)
}
