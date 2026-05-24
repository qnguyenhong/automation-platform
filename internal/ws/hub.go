package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	logger  *slog.Logger
}

// Client represents a single WebSocket connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		logger:  logger,
	}
}

// Run starts the hub.
func (h *Hub) Run() {
	h.logger.Info("websocket hub started")
}

// Upgrade upgrades an HTTP connection to WebSocket.
func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request) error {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register(client)

	go client.writePump()
	go client.readPump()

	return nil
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(message any) {
	data, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("failed to marshal broadcast message", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// BroadcastTo sends a message to a specific client.
func (h *Hub) BroadcastTo(client *Client, message any) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	select {
	case client.send <- data:
	default:
		close(client.send)
		delete(h.clients, client)
	}
}

func (h *Hub) register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true
	h.logger.Info("client connected", "total", len(h.clients))
}

func (h *Hub) unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		h.logger.Info("client disconnected", "total", len(h.clients))
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, _, err := c.conn.Read(nil)
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close(websocket.StatusNormalClosure, "")

	for message := range c.send {
		if err := c.conn.Write(nil, websocket.MessageText, message); err != nil {
			break
		}
	}
}
