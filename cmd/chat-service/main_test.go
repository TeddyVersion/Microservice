package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessageHandler_Success(t *testing.T) {
	msg := ChatMessage{From: "u1", To: "u2", Message: "hi"}
	body, _ := json.Marshal(msg)
	req := httptest.NewRequest(http.MethodPost, "/chat/send", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	sendMessageHandler(rec, req)
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

func TestSendMessageHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/chat/send", nil)
	rec := httptest.NewRecorder()
	sendMessageHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
