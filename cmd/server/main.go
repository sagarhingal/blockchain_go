package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"blockchain_go/internal/orderchain"
	"blockchain_go/internal/users"
	jwt "github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "secret"

// Server exposes HTTP handlers to interact with the permissioned chain.
type Server struct {
	chain  *orderchain.Chain
	mu     sync.Mutex
	users  *users.Store
	logger *log.Logger
}

func NewServer(dbPath ...string) *Server {
	path := "users.db"
	if len(dbPath) > 0 {
		path = dbPath[0]
	}
	store, err := users.NewStore(path)
	if err != nil {
		log.Fatalf("failed to init user store: %v", err)
	}
	// ensure root exists
	_ = store.CreateUser(users.User{Email: "root@example.com", Password: "12345", FirstName: "Root"})

	lf, _ := os.OpenFile("events.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	logger := log.New(lf, "", log.LstdFlags)

	return &Server{chain: orderchain.NewChain(), users: store, logger: logger}
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r)
	}
}

func (s *Server) tokenForUser(email string) (string, error) {
	claims := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(jwtSecret))
}

func (s *Server) userFromRequest(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", false
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	sub, ok := claims["sub"].(string)
	return sub, ok
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	var body users.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.users.CreateUser(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, _ := s.tokenForUser(body.Email)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.users.ValidateUser(creds.Email, creds.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, _ := s.tokenForUser(creds.Email)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ord := s.chain.CreateOrder(user)
	s.logger.Printf("order %s created by %s", ord.ID, user)
	json.NewEncoder(w).Encode(ord)
}

func (s *Server) manageRoles(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[1] != "roles" {
		http.NotFound(w, r)
		return
	}
	ord, ok2 := s.chain.Get(parts[0])
	if !ok2 {
		http.NotFound(w, r)
		return
	}
	var body struct{ Actor, Role string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ord.AddRole(user, body.Actor, body.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Printf("order %s role %s added for %s", ord.ID, body.Role, body.Actor)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) inviteWatcher(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[1] != "invite" {
		http.NotFound(w, r)
		return
	}
	ord, ok2 := s.chain.Get(parts[0])
	if !ok2 {
		http.NotFound(w, r)
		return
	}
	var body struct{ Actor string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ord.AddWatcher(user, body.Actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Printf("order %s watcher %s added by %s", ord.ID, body.Actor, user)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) updateStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[1] != "status" {
		http.NotFound(w, r)
		return
	}
	ord, ok2 := s.chain.Get(parts[0])
	if !ok2 {
		http.NotFound(w, r)
		return
	}
	var body struct{ Status string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ord.UpdateStatus(user, body.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Printf("order %s status %s by %s", ord.ID, body.Status, user)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) addAddon(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[1] != "addon" {
		http.NotFound(w, r)
		return
	}
	ord, ok2 := s.chain.Get(parts[0])
	if !ok2 {
		http.NotFound(w, r)
		return
	}
	var body struct{ Details string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ord.AddAddon(user, body.Details); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.logger.Printf("order %s addon by %s", ord.ID, user)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) listActors(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.userFromRequest(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := s.users.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rows)
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := body.Email
	if u, ok := s.userFromRequest(r); ok && email == "" {
		email = u
	}
	if email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}
	if err := s.users.UpdatePassword(email, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) getEvents(w http.ResponseWriter, r *http.Request) {
	user, ok := s.userFromRequest(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/order/")
	parts := strings.SplitN(id, "/", 2)
	if len(parts) < 2 || parts[1] != "events" {
		http.NotFound(w, r)
		return
	}
	ord, ok2 := s.chain.Get(parts[0])
	if !ok2 {
		http.NotFound(w, r)
		return
	}
	ev, err := ord.GetEvents(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ev)
}

func (s *Server) getChain(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.userFromRequest(r); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	json.NewEncoder(w).Encode(s.chain.Orders())
}

func main() {
	srv := NewServer()

	http.HandleFunc("/signup", withCORS(srv.signup))
	http.HandleFunc("/login", withCORS(srv.login))
	http.HandleFunc("/reset", withCORS(srv.resetPassword))
	http.HandleFunc("/order", withCORS(srv.createOrder))
	http.HandleFunc("/order/", withCORS(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	http.HandleFunc("/chain", withCORS(srv.getChain))
	http.HandleFunc("/transaction", withCORS(srv.createOrder))
	http.HandleFunc("/marketplace", withCORS(srv.listActors))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
