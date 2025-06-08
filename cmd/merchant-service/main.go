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

func payMerchantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MerchantPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// TODO: Add merchant payment logic here
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"payment_id": "merch123"}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/merchant/pay", payMerchantHandler)
	fmt.Println("merchant-service running on :8086")
	http.ListenAndServe(":8086", nil)
}
