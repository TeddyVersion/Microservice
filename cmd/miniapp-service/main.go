package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MiniApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

var miniapps = []MiniApp{
	{ID: "m1", Name: "Taxi", URL: "https://taxi.example.com"},
	{ID: "m2", Name: "Food Delivery", URL: "https://food.example.com"},
}

type MiniAppAuthRequest struct {
	UserID  string `json:"user_id"`
	MiniApp string `json:"miniapp_id"`
}

type MiniAppPaymentRequest struct {
	UserID    string  `json:"user_id"`
	MiniAppID string  `json:"miniapp_id"`
	Amount    float64 `json:"amount"`
}

var miniappAuthorizations = []MiniAppAuthRequest{}
var miniappPayments = []MiniAppPaymentRequest{}

func listMiniAppsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: miniapps})
}

func authorizeMiniAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MiniAppAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	miniappAuthorizations = append(miniappAuthorizations, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func payMiniAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MiniAppPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	miniappPayments = append(miniappPayments, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listMiniAppPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	miniAppID := r.URL.Query().Get("miniapp_id")
	var payments []MiniAppPaymentRequest
	for _, p := range miniappPayments {
		if p.MiniAppID == miniAppID {
			payments = append(payments, p)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: payments})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/miniapps", listMiniAppsHandler)
	http.HandleFunc("/miniapp/authorize", authorizeMiniAppHandler)
	http.HandleFunc("/miniapp/pay", payMiniAppHandler)
	http.HandleFunc("/miniapp/payments", listMiniAppPaymentsHandler)
	fmt.Println("miniapp-service running on :8090")
	http.ListenAndServe(":8090", nil)
}
