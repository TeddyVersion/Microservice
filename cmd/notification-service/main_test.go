package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendNotificationHandler(t *testing.T) {
	notifications = []Notification{} // reset
	body := bytes.NewReader([]byte(`{"user_id":"u1","type":"email","message":"test"}`))
	req := httptest.NewRequest(http.MethodPost, "/notify", body)
	rec := httptest.NewRecorder()
	sendNotificationHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}

func TestListNotificationsHandler(t *testing.T) {
	notifications = []Notification{{ID: "notif1", UserID: "u1", Type: "email", Message: "test", Status: "sent"}}
	req := httptest.NewRequest(http.MethodGet, "/notifications?user_id=u1", nil)
	rec := httptest.NewRecorder()
	listNotificationsHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}

func TestResendNotificationHandler(t *testing.T) {
	notifications = []Notification{{ID: "notif1", UserID: "u1", Type: "email", Message: "test", Status: "sent"}}
	req := httptest.NewRequest(http.MethodPost, "/notification/resend?id=notif1", nil)
	rec := httptest.NewRecorder()
	resendNotificationHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}

func TestSendNotificationHandler_Success(t *testing.T) {
	body := strings.NewReader(`{"user_id":"u1","type":"email","message":"hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/notify", body)
	rec := httptest.NewRecorder()
	sendNotificationHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}
