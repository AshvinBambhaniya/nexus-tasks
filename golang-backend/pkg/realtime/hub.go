package realtime

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	// Map Topic -> Map Connection -> bool
	clients map[string]map[*websocket.Conn]bool

	// Register requests from the clients.
	register chan subscription

	// Unregister requests from clients.
	unregister chan subscription

	// Broadcast messages to the clients.
	broadcast chan message

	// Mutex to protect clients map (though channel architecture mitigates this,
	// reading count or map access outside loop requires lock)
	mu sync.RWMutex

	logger *zap.Logger
}

type subscription struct {
	topic string
	conn  *websocket.Conn
}

type message struct {
	topic   string
	payload interface{}
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan message),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		clients:    make(map[string]map[*websocket.Conn]bool),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[sub.topic]; !ok {
				h.clients[sub.topic] = make(map[*websocket.Conn]bool)
			}
			h.clients[sub.topic][sub.conn] = true
			h.logger.Debug("Client subscribed", zap.String("topic", sub.topic))
			h.mu.Unlock()

		case sub := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[sub.topic]; ok {
				if _, ok := h.clients[sub.topic][sub.conn]; ok {
					delete(h.clients[sub.topic], sub.conn)
					h.logger.Debug("Client unsubscribed", zap.String("topic", sub.topic))
					if len(h.clients[sub.topic]) == 0 {
						delete(h.clients, sub.topic)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients, ok := h.clients[msg.topic]
			h.mu.RUnlock()

			if ok {
				for client := range clients {
					if err := client.WriteJSON(msg.payload); err != nil {
						h.logger.Error("error writing json to websocket", zap.Error(err))
						client.Close()
						// We should schedule unregister, but careful of deadlock if channel is full.
						// Typically in this loop, we can just remove it directly if we had the lock,
						// but since we are iterating, we need to be safe.
						// For now, let's just log. The reader loop in controller usually detects close.
					}
				}
			}
		}
	}
}

// Public API

func (h *Hub) Subscribe(topic string, conn *websocket.Conn) {
	h.register <- subscription{topic: topic, conn: conn}
}

func (h *Hub) Unsubscribe(topic string, conn *websocket.Conn) {
	h.unregister <- subscription{topic: topic, conn: conn}
}

func (h *Hub) Broadcast(topic string, payload interface{}) {
	h.broadcast <- message{topic: topic, payload: payload}
}
