package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"blockchain_go/internal/blockchain"
)

// swaggerSpec is a minimal OpenAPI definition describing the HTTP endpoints
// provided by the server. It is served at /swagger.json.
const swaggerSpec = `{
  "openapi": "3.0.0",
  "info": {
    "title": "Blockchain API",
    "version": "1.0"
  },
  "paths": {
    "/transaction": {
      "post": {
        "summary": "Add transaction",
        "responses": {
          "201": {"description": "created"}
        }
      }
    },
    "/chain": {
      "get": {
        "summary": "Get blockchain",
        "responses": {"200": {"description": "chain"}}
      }
    },
    "/validate": {
      "get": {
        "summary": "Validate chain",
        "responses": {"200": {"description": "status"}}
      }
    }
  }
}`

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

// serveSwagger writes the OpenAPI specification used to document the API.
func (s *Server) serveSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(swaggerSpec))
}

func main() {
	srv := NewServer(2)

	http.HandleFunc("/transaction", srv.addTransaction)
	http.HandleFunc("/chain", srv.getChain)
	http.HandleFunc("/validate", srv.validate)
	http.HandleFunc("/swagger.json", srv.serveSwagger)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
