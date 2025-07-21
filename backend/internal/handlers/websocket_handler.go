package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

// Generic Socket Message Structure
type SocketMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// Client represents a connected WebSocket client
type SocketClient struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan SocketMessage
	UserID *int32   // Optional: for user-specific targeting
	Groups []string // Optional: for group targeting
}

// SocketHandler manages WebSocket connections and events
type SocketHandler struct {
	// Client management
	clients    map[string]*SocketClient
	clientsMux sync.RWMutex

	// Event handlers
	eventHandlers map[string][]func(client *SocketClient, data interface{})
	handlersMux   sync.RWMutex

	// Channels
	register   chan *SocketClient
	unregister chan *SocketClient
	broadcast  chan SocketMessage
	done       chan struct{}
}

// NewSocketHandler creates a new socket handler
func NewSocketHandler() *SocketHandler {
	h := &SocketHandler{
		clients:       make(map[string]*SocketClient),
		eventHandlers: make(map[string][]func(*SocketClient, interface{})),
		register:      make(chan *SocketClient, 256),
		unregister:    make(chan *SocketClient, 256),
		broadcast:     make(chan SocketMessage, 256),
		done:          make(chan struct{}),
	}

	// Start the hub
	go h.run()

	return h
}

// ServeSocket handles WebSocket connections on /notifications endpoint
func (h *SocketHandler) ServeSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // In production: restrict this!
	})
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &SocketClient{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan SocketMessage, 256),
	}

	h.register <- client

	// Start goroutines for this client
	go h.writePump(client)
	go h.readPump(client)
}

// run manages the socket hub
func (h *SocketHandler) run() {
	defer close(h.done)

	for {
		select {
		case client := <-h.register:
			h.clientsMux.Lock()
			h.clients[client.ID] = client
			h.clientsMux.Unlock()
			log.Printf("Socket client connected: %s. Total clients: %d", client.ID, len(h.clients))

			// Send welcome message
			h.sendToClient(client, SocketMessage{
				Event: "connected",
				Data:  map[string]interface{}{"clientId": client.ID},
			})

		case client := <-h.unregister:
			h.clientsMux.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
			}
			h.clientsMux.Unlock()
			log.Printf("Socket client disconnected: %s. Total clients: %d", client.ID, len(h.clients))

		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

// readPump handles incoming messages from a client
func (h *SocketHandler) readPump(client *SocketClient) {
	defer func() {
		h.unregister <- client
		client.Conn.Close(websocket.StatusInternalError, "")
	}()

	ctx := context.Background()
	for {
		_, messageBytes, err := client.Conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			log.Printf("Socket read error: %v", err)
			return
		}

		// Parse incoming message
		var msg SocketMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Error parsing socket message: %v", err)
			continue
		}

		// Handle the message
		h.handleClientMessage(client, msg)
	}
}

// writePump sends messages to a client
func (h *SocketHandler) writePump(client *SocketClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close(websocket.StatusInternalError, "")
	}()

	ctx := context.Background()
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.Close(websocket.StatusNormalClosure, "")
				return
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			if err := client.Conn.Write(ctx, websocket.MessageText, messageBytes); err != nil {
				log.Printf("Socket write error: %v", err)
				return
			}

		case <-ticker.C:
			// Send ping
			if err := client.Conn.Ping(ctx); err != nil {
				return
			}
		}
	}
}

// handleClientMessage processes messages from clients
func (h *SocketHandler) handleClientMessage(client *SocketClient, msg SocketMessage) {
	log.Printf("Received message from client %s: %s", client.ID, msg.Event)

	// Handle built-in events
	switch msg.Event {
	case "ping":
		h.sendToClient(client, SocketMessage{
			Event: "pong",
			Data:  "pong",
		})
		return

	case "subscribe_user":
		// Client wants to subscribe to user-specific events
		if userID, ok := msg.Data.(float64); ok {
			userIDInt := int32(userID)
			client.UserID = &userIDInt
			log.Printf("Client %s subscribed to user %d", client.ID, userIDInt)
		}
		return

	case "join_group":
		// Client wants to join a group
		if groupName, ok := msg.Data.(string); ok {
			h.addClientToGroup(client, groupName)
			log.Printf("Client %s joined group %s", client.ID, groupName)
		}
		return

	case "leave_group":
		// Client wants to leave a group
		if groupName, ok := msg.Data.(string); ok {
			h.removeClientFromGroup(client, groupName)
			log.Printf("Client %s left group %s", client.ID, groupName)
		}
		return
	}

	// Handle custom events via registered handlers
	h.triggerEventHandlers(msg.Event, client, msg.Data)
}

// Event Handler Management

// OnEvent registers an event handler
func (h *SocketHandler) OnEvent(event string, handler func(client *SocketClient, data interface{})) {
	h.handlersMux.Lock()
	defer h.handlersMux.Unlock()

	if h.eventHandlers[event] == nil {
		h.eventHandlers[event] = make([]func(*SocketClient, interface{}), 0)
	}
	h.eventHandlers[event] = append(h.eventHandlers[event], handler)
	log.Printf("Event handler registered for: %s", event)
}

// triggerEventHandlers calls all registered handlers for an event
func (h *SocketHandler) triggerEventHandlers(event string, client *SocketClient, data interface{}) {
	h.handlersMux.RLock()
	handlers, exists := h.eventHandlers[event]
	h.handlersMux.RUnlock()

	if !exists {
		log.Printf("No handlers registered for event: %s", event)
		return
	}

	for _, handler := range handlers {
		go func(h func(*SocketClient, interface{})) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Event handler panic for %s: %v", event, r)
				}
			}()
			h(client, data)
		}(handler)
	}
}

// Broadcasting Methods

// EmitToAll sends event to all connected clients
func (h *SocketHandler) EmitToAll(event string, data interface{}) {
	message := SocketMessage{
		Event: event,
		Data:  data,
	}

	select {
	case h.broadcast <- message:
	default:
		log.Println("Broadcast channel is full")
	}
}

// EmitToUser sends event to specific user
func (h *SocketHandler) EmitToUser(userID int32, event string, data interface{}) {
	message := SocketMessage{
		Event: event,
		Data:  data,
	}

	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	for _, client := range h.clients {
		if client.UserID != nil && *client.UserID == userID {
			h.sendToClient(client, message)
		}
	}
}

// EmitToGroup sends event to all clients in a group
func (h *SocketHandler) EmitToGroup(groupName string, event string, data interface{}) {
	message := SocketMessage{
		Event: event,
		Data:  data,
	}

	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	for _, client := range h.clients {
		if h.isClientInGroup(client, groupName) {
			h.sendToClient(client, message)
		}
	}
}

// EmitToClient sends event to specific client
func (h *SocketHandler) EmitToClient(clientID string, event string, data interface{}) {
	message := SocketMessage{
		Event: event,
		Data:  data,
	}

	h.clientsMux.RLock()
	client, exists := h.clients[clientID]
	h.clientsMux.RUnlock()

	if exists {
		h.sendToClient(client, message)
	}
}

// Helper Methods

func (h *SocketHandler) sendToClient(client *SocketClient, message SocketMessage) {
	select {
	case client.Send <- message:
	default:
		// Client's send channel is full, remove it
		h.clientsMux.Lock()
		if _, ok := h.clients[client.ID]; ok {
			delete(h.clients, client.ID)
			close(client.Send)
		}
		h.clientsMux.Unlock()
	}
}

func (h *SocketHandler) broadcastToAll(message SocketMessage) {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	for _, client := range h.clients {
		h.sendToClient(client, message)
	}
}

func (h *SocketHandler) addClientToGroup(client *SocketClient, groupName string) {
	for _, group := range client.Groups {
		if group == groupName {
			return // Already in group
		}
	}
	client.Groups = append(client.Groups, groupName)
}

func (h *SocketHandler) removeClientFromGroup(client *SocketClient, groupName string) {
	for i, group := range client.Groups {
		if group == groupName {
			client.Groups = append(client.Groups[:i], client.Groups[i+1:]...)
			return
		}
	}
}

func (h *SocketHandler) isClientInGroup(client *SocketClient, groupName string) bool {
	for _, group := range client.Groups {
		if group == groupName {
			return true
		}
	}
	return false
}

// GetConnectedClients returns information about connected clients
func (h *SocketHandler) GetConnectedClients() map[string]interface{} {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	result := make(map[string]interface{})
	for id, client := range h.clients {
		clientInfo := map[string]interface{}{
			"id":     id,
			"groups": client.Groups,
		}
		if client.UserID != nil {
			clientInfo["userId"] = *client.UserID
		}
		result[id] = clientInfo
	}
	return result
}

// Shutdown gracefully closes the socket handler
func (h *SocketHandler) Shutdown() {
	close(h.done)

	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()

	for _, client := range h.clients {
		close(client.Send)
		client.Conn.Close(websocket.StatusServiceRestart, "Server shutting down")
	}
}
