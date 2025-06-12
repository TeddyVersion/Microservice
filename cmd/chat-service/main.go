package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type MoneyRequest struct {
	ID     string  `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type MoneyTransfer struct {
	ID     string  `json:"id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

var messages = []ChatMessage{}
var moneyRequests = []MoneyRequest{}
var moneyTransfers = []MoneyTransfer{}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var msg ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	messages = append(messages, msg)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"message_id": fmt.Sprintf("msg%d", len(messages))}})
}

func listMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: messages})
}

func getConversationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	var convo []ChatMessage
	for _, m := range messages {
		if (m.From == from && m.To == to) || (m.From == to && m.To == from) {
			convo = append(convo, m)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: convo})
}

func sendMoneyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MoneyTransfer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("tx%d", len(moneyTransfers)+1)
	req.Status = "completed"
	moneyTransfers = append(moneyTransfers, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func requestMoneyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req MoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	req.ID = fmt.Sprintf("req%d", len(moneyRequests)+1)
	req.Status = "pending"
	moneyRequests = append(moneyRequests, req)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: req})
}

func listMoneyRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var reqs []MoneyRequest
	for _, req := range moneyRequests {
		if req.To == userID || req.From == userID {
			reqs = append(reqs, req)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: reqs})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/chat/send", sendMessageHandler)
	http.HandleFunc("/chat/list", listMessagesHandler)
	http.HandleFunc("/chat/conversation", getConversationHandler)
	http.HandleFunc("/chat/sendmoney", sendMoneyHandler)
	http.HandleFunc("/chat/requestmoney", requestMoneyHandler)
	http.HandleFunc("/chat/moneyrequests", listMoneyRequestsHandler)
	fmt.Println("chat-service running on :8083")
	http.ListenAndServe(":8083", nil)
}
