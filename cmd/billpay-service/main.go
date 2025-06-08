package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gosmart/pkg/db"
	"gosmart/pkg/jwt"
	"gosmart/pkg/validation"
)

var sqlDB *sql.DB

func initDB() {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		panic("PG_DSN environment variable not set")
	}
	database, err := db.NewPostgres(dsn)
	if err != nil {
		panic(err)
	}
	if err := database.Ping(); err != nil {
		panic(err)
	}
	sqlDB = database
}

type BillPayment struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Biller    string  `json:"biller"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

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
	var paymentID int
	err := sqlDB.QueryRow(
		"INSERT INTO bill_payments (user_id, biller, amount, reference) VALUES ($1, $2, $3, $4) RETURNING id",
		req.UserID, req.Biller, req.Amount, req.Ref,
	).Scan(&paymentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Failed to create bill payment: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]interface{}{"payment_id": paymentID}})
}

func listBillPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	rows, err := sqlDB.Query("SELECT id, user_id, biller, amount, reference FROM bill_payments")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Failed to fetch bill payments"})
		return
	}
	defer rows.Close()
	var payments []BillPayment
	for rows.Next() {
		var p BillPayment
		if err := rows.Scan(&p.ID, &p.UserID, &p.Biller, &p.Amount, &p.Reference); err == nil {
			payments = append(payments, p)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: payments})
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
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/billpay/pay", payBillHandler)
	mux.HandleFunc("/billpay/list", listBillPaymentsHandler)
	fmt.Println("billpay-service running on :8084")
	http.ListenAndServe(":8084", authMiddleware(mux))
}
