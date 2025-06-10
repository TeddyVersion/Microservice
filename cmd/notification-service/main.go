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

type Notification struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

var notifications = []Notification{}

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
	n := Notification{
		ID:      fmt.Sprintf("notif%d", len(notifications)+1),
		UserID:  req.UserID,
		Type:    req.Type,
		Message: req.Message,
		Status:  "sent",
	}
	notifications = append(notifications, n)
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: n})
}

func listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	userID := r.URL.Query().Get("user_id")
	var userNotifs []Notification
	for _, n := range notifications {
		if n.UserID == userID {
			userNotifs = append(userNotifs, n)
		}
	}
	json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: userNotifs})
}

func resendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	for i, n := range notifications {
		if n.ID == id {
			notifications[i].Status = "resent"
			json.NewEncoder(w).Encode(APIResponse{Status: "success", Data: notifications[i]})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(APIResponse{Status: "error", Message: "Notification not found"})
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/notify", sendNotificationHandler)
	http.HandleFunc("/notifications", listNotificationsHandler)
	http.HandleFunc("/notification/resend", resendNotificationHandler)
	fmt.Println("notification-service running on :8091")
	http.ListenAndServe(":8091", nil)
}
