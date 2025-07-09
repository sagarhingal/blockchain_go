package orderchain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"
)

// Order represents a supply chain order with scoped participants and encrypted logs.
type Order struct {
	ID       string
	Owner    string
	Actors   map[string]string // username->role
	Status   string
	Created  time.Time
	Events   []event
	AddOns   []string
	Watchers map[string]bool
	key      []byte
	mu       sync.Mutex
}

type event struct {
	Time    time.Time
	Actor   string
	Message string // encrypted
}

// Chain manages a set of orders.
type Chain struct {
	orders map[string]*Order
	mu     sync.Mutex
}

// Orders returns all orders in the chain.
func (c *Chain) Orders() []*Order {
	c.mu.Lock()
	defer c.mu.Unlock()
	list := make([]*Order, 0, len(c.orders))
	for _, o := range c.orders {
		list = append(list, o)
	}
	return list
}

// NewChain returns an initialized permissioned chain.
func NewChain() *Chain {
	return &Chain{orders: make(map[string]*Order)}
}

// CreateOrder creates a new order owned by username.
func (c *Chain) CreateOrder(owner string) *Order {
	id := randomID()
	key := make([]byte, 32)
	rand.Read(key)
	ord := &Order{
		ID:       id,
		Owner:    owner,
		Actors:   map[string]string{owner: "client"},
		Status:   "created",
		Created:  time.Now(),
		key:      key,
		Watchers: make(map[string]bool),
	}
	c.mu.Lock()
	c.orders[id] = ord
	c.mu.Unlock()
	ord.addEvent(owner, "order created")
	return ord
}

// Get returns the order by ID.
func (c *Chain) Get(id string) (*Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ord, ok := c.orders[id]
	return ord, ok
}

// AddRole assigns a role to an actor. Only the owner can manage roles.
func (o *Order) AddRole(owner, actor, role string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if owner != o.Owner {
		return errors.New("not order owner")
	}
	o.Actors[actor] = role
	o.addEvent(owner, "added role "+role+" for "+actor)
	return nil
}

// UpdateStatus sets a new status if actor has a role.
func (o *Order) UpdateStatus(actor, status string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, ok := o.Actors[actor]; !ok {
		return errors.New("unauthorized actor")
	}
	o.Status = status
	o.addEvent(actor, "status updated to "+status)
	return nil
}

// AddAddon records an add-on request from an actor.
func (o *Order) AddAddon(actor, details string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, ok := o.Actors[actor]; !ok {
		return errors.New("unauthorized actor")
	}
	o.AddOns = append(o.AddOns, details)
	o.addEvent(actor, "add-on: "+details)
	return nil
}

// AddWatcher registers a watcher who will receive event visibility.
func (o *Order) AddWatcher(owner, watcher string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if owner != o.Owner {
		return errors.New("not order owner")
	}
	o.Watchers[watcher] = true
	o.addEvent(owner, "added watcher "+watcher)
	return nil
}

// IsParticipant checks if user can interact with the order.
func (o *Order) IsParticipant(user string) bool {
	if user == o.Owner {
		return true
	}
	if _, ok := o.Actors[user]; ok {
		return true
	}
	return false
}

// GetEvents returns decrypted events.
func (o *Order) GetEvents(requester string) ([]map[string]string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.IsParticipant(requester) && !o.Watchers[requester] {
		return nil, errors.New("unauthorized")
	}
	events := make([]map[string]string, len(o.Events))
	for i, e := range o.Events {
		msg, err := decrypt(o.key, e.Message)
		if err != nil {
			return nil, err
		}
		events[i] = map[string]string{"time": e.Time.Format(time.RFC3339), "actor": e.Actor, "message": msg}
	}
	return events, nil
}

func (o *Order) addEvent(actor, msg string) {
	enc, _ := encrypt(o.key, msg)
	o.Events = append(o.Events, event{Time: time.Now(), Actor: actor, Message: enc})
}

func randomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func encrypt(key []byte, msg string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, enc string) (string, error) {
	data, err := hex.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return string(data), nil
}
