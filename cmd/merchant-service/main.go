package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MerchantPaymentRequest struct {
	UserID     string  `json:"user_id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Method     string  `json:"method"` // QR or TILL
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type MerchantPayment struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Method     string  `json:"method"` // QR or TILL
	Status     string  `json:"status"`
}

var merchantPayments = []MerchantPayment{}

func payMerchantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MerchantPayment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("mp%d", len(merchantPayments)+1)
	req.Status = "completed"
	merchantPayments = append(merchantPayments, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listMerchantPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userPayments []MerchantPayment
	for _, p := range merchantPayments {
		if p.UserID == userID {
			userPayments = append(userPayments, p)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userPayments})
}

func getMerchantStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	merchantID := r.URL.Query().Get("merchant_id")
	var total float64
	var count int
	for _, p := range merchantPayments {
		if p.MerchantID == merchantID {
			total += p.Amount
			count++
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]interface{}{"merchant_id": merchantID, "total": total, "count": count}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/merchant/pay", payMerchantHandler)
	http.HandleFunc("/merchant/list", listMerchantPaymentsHandler)
	http.HandleFunc("/merchant/stats", getMerchantStatsHandler)
	fmt.Println("merchant-service running on :8086")
	http.ListenAndServe(":8086", nil)
}
