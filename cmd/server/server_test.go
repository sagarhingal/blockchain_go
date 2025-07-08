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
	mux.HandleFunc("/transaction", srv.addTransaction)
	mux.HandleFunc("/chain", srv.getChain)
	mux.HandleFunc("/validate", srv.validate)
	mux.HandleFunc("/swagger.json", srv.serveSwagger)

	// Add a transaction
	tx := map[string]interface{}{"from": "A", "to": "B", "amount": 1.0}
	body, _ := json.Marshal(tx)
	req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// Fetch the chain
	req = httptest.NewRequest(http.MethodGet, "/chain", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Validate chain
	req = httptest.NewRequest(http.MethodGet, "/validate", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Fetch swagger spec
	req = httptest.NewRequest(http.MethodGet, "/swagger.json", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("openapi")) {
		t.Fatal("swagger spec missing 'openapi' key")
	}
}
