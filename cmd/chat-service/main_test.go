package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessageHandler(t *testing.T) {
	messages = []ChatMessage{} // reset
	body := bytes.NewReader([]byte(`{"from":"u1","to":"u2","message":"hi"}`))
	req := httptest.NewRequest(http.MethodPost, "/chat/send", body)
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

func TestListMessagesHandler(t *testing.T) {
	messages = []ChatMessage{{From: "u1", To: "u2", Message: "hi"}}
	req := httptest.NewRequest(http.MethodGet, "/chat/list", nil)
	rec := httptest.NewRecorder()
	listMessagesHandler(rec, req)
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

func TestGetConversationHandler(t *testing.T) {
	messages = []ChatMessage{{From: "u1", To: "u2", Message: "hi"}, {From: "u2", To: "u1", Message: "hello"}, {From: "u3", To: "u1", Message: "other"}}
	req := httptest.NewRequest(http.MethodGet, "/chat/conversation?from=u1&to=u2", nil)
	rec := httptest.NewRecorder()
	getConversationHandler(rec, req)
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
	// Should return 2 messages for the conversation
	if len(resp.Data.([]interface{})) != 2 {
		t.Errorf("expected 2 messages, got %d", len(resp.Data.([]interface{})))
	}
}
