package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransferHandler_Success(t *testing.T) {
	body := strings.NewReader(`{"sender_id":"u1","recipient_id":"u2","amount":100,"type":"CBE"}`)
	req := httptest.NewRequest(http.MethodPost, "/transfer", body)
	rec := httptest.NewRecorder()
	transferHandler(rec, req)
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
