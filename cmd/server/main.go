package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"blockchain_go/internal/blockchain"
)

// Server exposes HTTP handlers to interact with the blockchain.

type Server struct {
	chain blockchain.Blockchain
	mu    sync.Mutex
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
	s.mu.Lock()
	s.chain.AddBlock(tx.From, tx.To, tx.Amount)
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// getChain responds with the current blockchain.
func (s *Server) getChain(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]bool{"valid": valid}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	srv := NewServer(2)

	http.HandleFunc("/transaction", srv.addTransaction)
	http.HandleFunc("/chain", srv.getChain)
	http.HandleFunc("/validate", srv.validate)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
