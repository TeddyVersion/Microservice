package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"` // email, sms, in-app
	Message string `json:"message"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func sendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Invalid request"})
		return
	}
	// TODO: Add notification logic here
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: map[string]string{"notification_id": "notif123"}})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/notify", sendNotificationHandler)
	fmt.Println("notification-service running on :8091")
	http.ListenAndServe(":8091", nil)
}
