package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()
	getAccountsHandler(rec, req)
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

func TestCreateAccountHandler(t *testing.T) {
	accounts = []Account{} // reset
	body := bytes.NewReader([]byte(`{"balance": 500, "type": "savings"}`))
	req := httptest.NewRequest(http.MethodPost, "/accounts", body)
	rec := httptest.NewRecorder()
	createAccountHandler(rec, req)
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

func TestGetAccountByIDHandler(t *testing.T) {
	accounts = []Account{{ID: "acc1", Balance: 100, Type: "savings"}}
	req := httptest.NewRequest(http.MethodGet, "/account?id=acc1", nil)
	rec := httptest.NewRecorder()
	getAccountByIDHandler(rec, req)
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
