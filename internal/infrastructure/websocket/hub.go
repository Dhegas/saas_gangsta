package websocket

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins since middleware and JWT authentication guard it
	},
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool
	logger  *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[string]map[*websocket.Conn]bool),
		logger:  logger,
	}
}

func (h *Hub) Register(clientID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.clients[clientID]; !exists {
		h.clients[clientID] = make(map[*websocket.Conn]bool)
	}
	h.clients[clientID][conn] = true
	h.logger.Info("WebSocket client registered", "clientID", clientID)
}

func (h *Hub) Unregister(clientID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, exists := h.clients[clientID]; exists {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, clientID)
		}
	}
	h.logger.Info("WebSocket client unregistered", "clientID", clientID)
}

func (h *Hub) SendTo(clientID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns, exists := h.clients[clientID]
	if !exists {
		h.logger.Warn("WebSocket SendTo: client not online", "clientID", clientID)
		return
	}
	for conn := range conns {
		err := conn.WriteJSON(message)
		if err != nil {
			h.logger.Error("WebSocket SendTo: failed to write JSON to connection", "clientID", clientID, "error", err)
		}
	}
}
