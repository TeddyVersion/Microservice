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

var messages = []ChatMessage{}

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

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/chat/send", sendMessageHandler)
	http.HandleFunc("/chat/list", listMessagesHandler)
	http.HandleFunc("/chat/conversation", getConversationHandler)
	fmt.Println("chat-service running on :8083")
	http.ListenAndServe(":8083", nil)
}
