package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"blockchain_go/internal/users"
)

func TestOrderFlow(t *testing.T) {
	srv := NewServer(":memory:")
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", srv.signup)
	mux.HandleFunc("/login", srv.login)
	mux.HandleFunc("/order", srv.createOrder)
	mux.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/roles") {
			srv.manageRoles(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/status") {
			srv.updateStatus(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/addon") {
			srv.addAddon(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/events") {
			srv.getEvents(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/invite") {
			srv.inviteWatcher(w, r)
			return
		}
		http.NotFound(w, r)
	})
	mux.HandleFunc("/chain", srv.getChain)

	signup := func(u users.User) string {
		body, _ := json.Marshal(u)
		req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("signup failed: %d %s", w.Code, w.Body.String())
		}
		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		return resp["token"]
	}

	aliceTok := signup(users.User{Email: "alice@example.com", Password: "pass", FirstName: "A"})
	bobTok := signup(users.User{Email: "bob@example.com", Password: "pass", FirstName: "B"})
	charTok := signup(users.User{Email: "char@example.com", Password: "pass", FirstName: "C"})

	// alice creates order
	req := httptest.NewRequest(http.MethodPost, "/order", nil)
	req.Header.Set("Authorization", "Bearer "+aliceTok)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("order create failed: %d %s", w.Code, w.Body.String())
	}
	var ord map[string]interface{}
	json.NewDecoder(w.Body).Decode(&ord)
	id := ord["ID"].(string)

	// alice adds bob role supplier
	body, _ := json.Marshal(map[string]string{"Actor": "bob@example.com", "Role": "supplier"})
	req = httptest.NewRequest(http.MethodPost, "/order/"+id+"/roles", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+aliceTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("add role failed: %d %s", w.Code, w.Body.String())
	}

	// alice invites char as watcher
	body, _ = json.Marshal(map[string]string{"Actor": "char@example.com"})
	req = httptest.NewRequest(http.MethodPost, "/order/"+id+"/invite", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+aliceTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("invite failed: %d %s", w.Code, w.Body.String())
	}

	// bob updates status
	body, _ = json.Marshal(map[string]string{"Status": "shipped"})
	req = httptest.NewRequest(http.MethodPost, "/order/"+id+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+bobTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status update failed: %d %s", w.Code, w.Body.String())
	}

	// char attempts status update - expect unauthorized
	req = httptest.NewRequest(http.MethodPost, "/order/"+id+"/status", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+charTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Fatalf("char should not update status")
	}

	// char fetches events
	req = httptest.NewRequest(http.MethodGet, "/order/"+id+"/events", nil)
	req.Header.Set("Authorization", "Bearer "+charTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("char cannot get events: %d", w.Code)
	}

	// fetch chain
	req = httptest.NewRequest(http.MethodGet, "/chain", nil)
	req.Header.Set("Authorization", "Bearer "+aliceTok)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("chain fetch failed: %d", w.Code)
	}
}
