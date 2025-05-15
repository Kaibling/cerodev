package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type websocketRepo interface {
	Add(token string, conn *websocket.Conn)
	RemoveAndClose(token string)
	SendJSON(data any, token string)
	HealthCheckAll()
}

type WebSocketService struct {
	repo websocketRepo
}

func NewWebSocketService(repo websocketRepo) *WebSocketService {
	return &WebSocketService{repo: repo}
}

func (s *WebSocketService) Add(token string, conn *websocket.Conn) {
	s.repo.Add(token, conn)
}

func (s *WebSocketService) RemoveAndClose(token string) {
	s.repo.RemoveAndClose(token)
}

func (s *WebSocketService) SendJSON(data any, token string) {
	s.repo.SendJSON(data, token)
}

func (s *WebSocketService) StartHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick at", time.Now())
			s.repo.HealthCheckAll()
		case <-ctx.Done():
			fmt.Println("stopping web socket")
			return
		}
	}
}
