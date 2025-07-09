package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"blockchain_go/internal/blockchain"
	"blockchain_go/internal/users"
)

// Server exposes HTTP handlers to interact with the blockchain.

type Server struct {
	chain    blockchain.Blockchain
	mu       sync.Mutex
	users    *users.Store
	sessions map[string]string
	sessMu   sync.Mutex
}

// withCORS wraps an http.HandlerFunc adding permissive CORS headers so the
// web UI running on a different port can interact with the API.
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r)
	}
}

// withLogging measures the time taken to execute the handler and logs the
// request method, path and duration.
func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start)
		log.Printf("%s %s took %s", r.Method, r.URL.Path, duration)
	}
}

// NewServer returns a Server with a freshly created blockchain using the
// provided mining difficulty.
func NewServer(difficulty int) *Server {
	store, err := users.NewStore("users.db")
	if err != nil {
		log.Fatalf("failed to init user store: %v", err)
	}
	// ensure root user exists
	_ = store.CreateUser("root", "12345")
	return &Server{
		chain:    blockchain.CreateBlockchain(difficulty),
		users:    store,
		sessions: make(map[string]string),
	}
}

func (s *Server) usernameFromRequest(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}
	s.sessMu.Lock()
	defer s.sessMu.Unlock()
	u, ok := s.sessions[cookie.Value]
	return u, ok
}

func (s *Server) newSession(username string) string {
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)
	s.sessMu.Lock()
	s.sessions[token] = username
	s.sessMu.Unlock()
	return token
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.users.ValidateUser(creds.Username, creds.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token := s.newSession(creds.Username)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, Path: "/"})
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.users.CreateUser(creds.Username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := s.newSession(creds.Username)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, Path: "/"})
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	user, ok := s.usernameFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.users.UpdatePassword(user, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		s.sessMu.Lock()
		delete(s.sessions, cookie.Value)
		s.sessMu.Unlock()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// addTransaction handles POST requests creating a transaction block.
func (s *Server) addTransaction(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.usernameFromRequest(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var tx struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	log.Printf("New transaction from %s to %s amount %.2f", tx.From, tx.To, tx.Amount)
	s.mu.Lock()
	s.chain.AddBlock(tx.From, tx.To, tx.Amount)
	log.Printf("Block added. Chain length %d", len(s.chain.Chain))
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// getChain responds with the current blockchain.
func (s *Server) getChain(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.usernameFromRequest(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("Chain requested. Length %d", len(s.chain.Chain))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.chain); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validate checks that the blockchain state is valid.
func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.usernameFromRequest(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	s.mu.Lock()
	valid := s.chain.IsValid()
	s.mu.Unlock()
	log.Printf("Chain validation result: %v", valid)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"valid": valid}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	srv := NewServer(2)

	http.HandleFunc("/login", withCORS(withLogging(srv.login)))
	http.HandleFunc("/signup", withCORS(withLogging(srv.signup)))
	http.HandleFunc("/reset", withCORS(withLogging(srv.resetPassword)))
	http.HandleFunc("/logout", withCORS(withLogging(srv.logout)))
	http.HandleFunc("/transaction", withCORS(withLogging(srv.addTransaction)))
	http.HandleFunc("/chain", withCORS(withLogging(srv.getChain)))
	http.HandleFunc("/validate", withCORS(withLogging(srv.validate)))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
