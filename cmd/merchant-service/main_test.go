package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPayMerchantHandler_Success(t *testing.T) {
	body := strings.NewReader(`{"user_id":"u1","merchant_id":"m1","amount":50,"method":"QR"}`)
	req := httptest.NewRequest(http.MethodPost, "/merchant/pay", body)
	rec := httptest.NewRecorder()
	payMerchantHandler(rec, req)
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
