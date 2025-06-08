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

	"golang.org/x/crypto/bcrypt"
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

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// In-memory store for demonstration
var users = []RegisterRequest{}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	if err := validation.RequireFields(map[string]string{"phone": req.Phone, "password": req.Password}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Failed to hash password"})
		return
	}
	var userID int
	err = sqlDB.QueryRow(
		"INSERT INTO users (phone, password_hash) VALUES ($1, $2) RETURNING id",
		req.Phone, string(hash),
	).Scan(&userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Failed to create user: " + err.Error()})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]interface{}{"user_id": userID}})
}

type User struct {
	ID       int    `json:"id"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	rows, err := sqlDB.Query("SELECT id, phone FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Failed to fetch users"})
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Phone); err == nil {
			users = append(users, u)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: users})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/register" || r.URL.Path == "/healthz" {
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
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/users", listUsersHandler)
	fmt.Println("auth-service running on :8080")
	http.ListenAndServe(":8080", authMiddleware(mux))
}
