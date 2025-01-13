package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"wasa/service/database"

	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	userID   int
	send     chan []byte
	hub      *Hub
	isActive bool
	mu       sync.RWMutex
}

type Hub struct {
	clients             map[*Client]bool
	broadcast           chan []byte
	register            chan *Client
	unregister          chan *Client
	userConnections     map[int][]*Client
	conversationClients map[int][]*Client
	mu                  sync.RWMutex
	db                  database.AppDatabase
	reconnectTimeout    time.Duration
}

func NewHub(db database.AppDatabase) *Hub {
	return &Hub{
		clients:             make(map[*Client]bool),
		broadcast:           make(chan []byte),
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		userConnections:     make(map[int][]*Client),
		conversationClients: make(map[int][]*Client),
		db:                  db,
		reconnectTimeout:    5 * time.Minute,
	}
}

type WebSocketMessage struct {
	Type           string                 `json:"type"`
	ConversationID int                    `json:"conversation_id"`
	Payload        map[string]interface{} `json:"payload"`
	Timestamp      time.Time              `json:"timestamp"`
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.userConnections[client.userID] = append(h.userConnections[client.userID], client)
			h.mu.Unlock()

			h.AddClientToConversations(client)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				connections := h.userConnections[client.userID]
				for i, conn := range connections {
					if conn == client {
						h.userConnections[client.userID] = append(connections[:i], connections[i+1:]...)
						break
					}
				}

				for convID, clients := range h.conversationClients {
					for i, c := range clients {
						if c == client {
							h.conversationClients[convID] = append(clients[:i], clients[i+1:]...)
							break
						}
					}
					if len(h.conversationClients[convID]) == 0 {
						delete(h.conversationClients, convID)
					}
				}

				if len(h.userConnections[client.userID]) == 0 {
					delete(h.userConnections, client.userID)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID int, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.userConnections[userID]; ok {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) SendToConversation(conversationID int, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	fmt.Printf("Attempting to send message to conversation %d, %d clients\n",
		conversationID, len(h.conversationClients[conversationID]))

	if clients, ok := h.conversationClients[conversationID]; ok {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) AddClientToConversations(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conversations, err := h.db.GetAllConversations(client.userID)
	if err != nil {
		fmt.Printf("Error getting conversations for user %d: %v\n", client.userID, err)
		return
	}

	for _, conv := range conversations {
		convID := conv.ID
		exists := false
		for _, existingClient := range h.conversationClients[convID] {
			if existingClient == client {
				exists = true
				break
			}
		}

		if !exists {
			h.conversationClients[convID] = append(h.conversationClients[convID], client)
			fmt.Printf("Added client to conversation %d, total clients: %d\n",
				convID, len(h.conversationClients[convID]))
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		fmt.Printf("Message received in readPump: %s\n", string(message))

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			continue
		}

		switch wsMessage.Type {
		case "message":
			c.hub.broadcast <- message
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (rt *APIRouter) HandleWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("WebSocket connection attempt received")

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.Header.Set("Authorization", "Bearer "+token)
	userID, err := getToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to upgrade connection")
		return
	}

	client := &Client{
		conn:     conn,
		userID:   userID,
		hub:      rt.wsHub,
		send:     make(chan []byte, 256),
		isActive: true,
	}

	client.hub.register <- client

	go client.startHeartbeat()

	go client.readPump()
	go client.writePump()
}

func (c *Client) startHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		c.mu.RLock()
		if !c.isActive {
			c.mu.RUnlock()
			return
		}
		c.mu.RUnlock()

		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
		<-ticker.C
	}
}
