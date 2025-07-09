package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"blockchain_go/internal/users"
)

func TestOrderFlow(t *testing.T) {
	srv := NewServer(":memory:")
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", srv.signup)
	mux.HandleFunc("/login", srv.login)
	mux.HandleFunc("/order", srv.createOrder)
	mux.HandleFunc("/chain", srv.getChain)

	// signup new user
	user := users.User{Username: "alice" + fmt.Sprint(time.Now().UnixNano()), Password: "pass", FirstName: "Alice"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("signup failed: %d %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	token := resp["token"]

	// create order
	req = httptest.NewRequest(http.MethodPost, "/order", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("order create failed: %d %s", w.Code, w.Body.String())
	}

	// fetch chain
	req = httptest.NewRequest(http.MethodGet, "/chain", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("chain fetch failed: %d", w.Code)
	}
}
