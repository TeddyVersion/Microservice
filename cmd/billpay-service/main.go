package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gosmart/pkg/jwt"
	"gosmart/pkg/validation"
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
	if err := validation.RequireFields(map[string]string{"user_id": req.UserID, "biller": req.Biller, "reference": req.Ref}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: err.Error()})
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Missing token"})
			return
		}
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		_, err := jwt.ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid token"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/billpay/pay", payBillHandler)
	mux.HandleFunc("/billpay/list", listBillPaymentsHandler)
	fmt.Println("billpay-service running on :8084")
	http.ListenAndServe(":8084", authMiddleware(mux))
}
