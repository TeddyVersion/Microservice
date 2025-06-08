package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPayBillHandler_Success(t *testing.T) {
	billPayments = []BillPaymentRequest{} // reset in-memory store
	requestBody := BillPaymentRequest{
		UserID: "u1",
		Biller: "DSTV",
		Amount: 100.0,
		Ref:    "INV123",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/billpay/pay", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	payBillHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
}

func TestPayBillHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/billpay/pay", nil)
	rec := httptest.NewRecorder()
	payBillHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestListBillPaymentsHandler(t *testing.T) {
	billPayments = []BillPaymentRequest{{UserID: "u1", Biller: "DSTV", Amount: 100, Ref: "INV123"}}
	req := httptest.NewRequest(http.MethodGet, "/billpay/list", nil)
	rec := httptest.NewRecorder()
	listBillPaymentsHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var resp APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}
	if len(resp.Data.([]interface{})) == 0 {
		t.Errorf("expected at least one bill payment")
	}
}
