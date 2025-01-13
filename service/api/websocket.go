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

	clients := h.conversationClients[conversationID]
	fmt.Printf("Sending message to conversation %d (%d clients)\n",
		conversationID, len(clients))

	for _, client := range clients {
		select {
		case client.send <- message:
			fmt.Printf("Sent message to client %d\n", client.userID)
		default:
			fmt.Printf("Failed to send message to client %d, closing connection\n", client.userID)
			close(client.send)
			delete(h.clients, client)
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

func (h *Hub) AddConversationClient(conversationID int, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, existingClient := range h.conversationClients[conversationID] {
		if existingClient == client {
			return
		}
	}

	h.conversationClients[conversationID] = append(h.conversationClients[conversationID], client)
	fmt.Printf("Added client to conversation %d, total clients: %d\n",
		conversationID, len(h.conversationClients[conversationID]))
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

		fmt.Printf("Received message: %s\n", string(message))

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		switch wsMessage.Type {
		case "message":
			c.hub.broadcast <- message

		// todo: probably will need to update the database, if I am not updating
		// already, the entry in the messageSeen of the user sending it to waiting.
		// After each user sees it, they should update the table to seen.
		// However, I am not sure what entry they should be updating.
		// Because, initially, the sender should know the other host hasn't
		// read the message, and after the other host reads it, the sender should
		// know and update the value.
		// Will need to ask Claude this.
		// Basically, the sender should know the other host has read the message.
		// How we achieve this is the question: I know we need to broadcast the message,
		// but I am thinking more along the line of persistency in the db.

		case "messages_seen":
			fmt.Printf("Processing messages_seen: %+v\n", wsMessage)

			messageIDs, ok := wsMessage.Payload["message_ids"].([]interface{})
			if !ok {
				fmt.Printf("Invalid message_ids format: %v\n", wsMessage.Payload["message_ids"])
				continue
			}

			ids := make([]int, 0, len(messageIDs))
			for _, v := range messageIDs {
				if id, ok := v.(float64); ok {
					ids = append(ids, int(id))
				}
			}

			if len(ids) == 0 {
				fmt.Printf("No valid message IDs found\n")
				continue
			}

			fmt.Printf("Marking messages %v as seen by user %d\n", ids, c.userID)

			err := c.hub.db.MarkMessagesSeen(c.userID, ids)
			if err != nil {
				fmt.Printf("Error marking messages as seen: %v\n", err)
				continue
			}

			wsMessage.Payload["user_id"] = c.userID
			updatedMessage, err := json.Marshal(wsMessage)
			if err != nil {
				fmt.Printf("Error marshaling updated message: %v\n", err)
				continue
			}

			fmt.Printf("Broadcasting seen status to conversation %d\n", wsMessage.ConversationID)
			c.hub.SendToConversation(wsMessage.ConversationID, updatedMessage)
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

func (c *Client) logState() {
	fmt.Printf("Client state - UserID: %d, Active: %v\n", c.userID, c.isActive)
	if c.hub != nil {
		fmt.Printf("Hub state - Connected clients: %d, Client's conversations: %d\n",
			len(c.hub.clients), len(c.hub.conversationClients))
	}
}
