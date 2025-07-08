package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"blockchain_go/internal/blockchain"
)

// Server exposes HTTP handlers to interact with the blockchain.

type Server struct {
	chain blockchain.Blockchain
	mu    sync.Mutex
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
	return &Server{chain: blockchain.CreateBlockchain(difficulty)}
}

// addTransaction handles POST requests creating a transaction block.
func (s *Server) addTransaction(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/transaction", withCORS(withLogging(srv.addTransaction)))
	http.HandleFunc("/chain", withCORS(withLogging(srv.getChain)))
	http.HandleFunc("/validate", withCORS(withLogging(srv.validate)))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
