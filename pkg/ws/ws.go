package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketRepo struct {
	mu      sync.RWMutex
	clients map[string]*websocket.Conn
}

func New() *WebSocketRepo {
	return &WebSocketRepo{clients: map[string]*websocket.Conn{}}
}

func (r *WebSocketRepo) Add(token string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[token] = conn
}

func (r *WebSocketRepo) RemoveAndClose(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if client, ok := r.clients[token]; ok {
		client.Close()
		delete(r.clients, token)
		fmt.Printf("Cleaned up client")
	}
}

func (r *WebSocketRepo) SendJSON(data any, token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[token].WriteJSON(data)
}

func (r *WebSocketRepo) HealthCheckAll() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Println("check sockets")
	for token, client := range r.clients {
		fmt.Printf("check socket %s", token)
		err := client.WriteControl(
			websocket.PingMessage,
			[]byte("ping"),
			time.Now().Add(10*time.Second),
		)
		if err != nil {
			fmt.Printf("Client dead or unresponsive: %v", err)
			go r.RemoveAndClose(token) // clean it up
		}
	}
}
