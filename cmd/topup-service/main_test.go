package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTopupHandler_Success(t *testing.T) {
	body := strings.NewReader(`{"user_id":"u1","amount":10,"type":"airtime","phone":"251900000000"}`)
	req := httptest.NewRequest(http.MethodPost, "/topup", body)
	rec := httptest.NewRecorder()
	topupHandler(rec, req)
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
