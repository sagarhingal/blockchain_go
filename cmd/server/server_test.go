package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerEndpoints(t *testing.T) {
	srv := NewServer(1)
	mux := http.NewServeMux()
	mux.HandleFunc("/login", srv.login)
	mux.HandleFunc("/transaction", srv.addTransaction)
	mux.HandleFunc("/chain", srv.getChain)
	mux.HandleFunc("/validate", srv.validate)

	// login as root
	creds := map[string]string{"username": "root", "password": "12345"}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: %d", w.Code)
	}
	cookie := w.Result().Cookies()[0]

	// Add a transaction
	tx := map[string]interface{}{"from": "A", "to": "B", "amount": 1.0}
	body, _ = json.Marshal(tx)
	req = httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewReader(body))
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Fetch the chain
	req = httptest.NewRequest(http.MethodGet, "/chain", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Validate chain
	req = httptest.NewRequest(http.MethodGet, "/validate", nil)
	req.AddCookie(cookie)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
