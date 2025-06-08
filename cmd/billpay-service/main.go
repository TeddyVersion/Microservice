package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BillPaymentRequest struct {
	UserID string  `json:"user_id"`
	Biller string  `json:"biller"`
	Amount float64 `json:"amount"`
	Ref    string  `json:"reference"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// In-memory store for demonstration
var billPayments = []BillPaymentRequest{}

func payBillHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req BillPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// Simulate storing the bill payment
	billPayments = append(billPayments, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"payment_id": fmt.Sprintf("pay%d", len(billPayments))}})
}

func listBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: billPayments})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/billpay/pay", payBillHandler)
	http.HandleFunc("/billpay/list", listBillPaymentsHandler)
	fmt.Println("billpay-service running on :8084")
	http.ListenAndServe(":8084", nil)
}
