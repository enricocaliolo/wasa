package api

import (
	"encoding/json"
	"sync"
	"time"
	"wasa/service/database"

	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
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
	logger              *logrus.Logger
}

func NewHub(db database.AppDatabase, logger *logrus.Logger) *Hub {
	return &Hub{
		clients:             make(map[*Client]bool),
		broadcast:           make(chan []byte),
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		userConnections:     make(map[int][]*Client),
		conversationClients: make(map[int][]*Client),
		db:                  db,
		reconnectTimeout:    5 * time.Minute,
		logger:              logger,
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
				client.mu.Lock()
				client.isActive = false
				client.mu.Unlock()

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

	for _, client := range clients {
		select {
		case client.send <- message:
		default:
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
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			c.hub.logger.WithField("recover", r).Error("Recovered in readPump")
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		c.hub.logger.WithError(err).Error("Error setting read deadline")
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			c.hub.logger.WithError(err).Error("Error setting read deadline")
			return err
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if !websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
					c.hub.logger.WithError(err).Error("Unexpected close error")
				}
			}
			break
		}

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			c.hub.logger.WithError(err).Error("Error unmarshaling message")
			continue
		}

		switch wsMessage.Type {
		case "message":
			c.hub.broadcast <- message

		case "messages_seen":
			messageIDs, ok := wsMessage.Payload["message_ids"].([]interface{})
			if !ok {
				c.hub.logger.WithError(err).Error("Invalid message_ids format: %w", wsMessage.Payload["message_ids"])
				continue
			}

			ids := make([]int, 0, len(messageIDs))
			for _, v := range messageIDs {
				if id, ok := v.(float64); ok {
					ids = append(ids, int(id))
				}
			}

			if len(ids) == 0 {
				c.hub.logger.WithError(err).Error("No message IDs found")
				continue
			}

			err := c.hub.db.MarkMessagesSeen(c.userID, ids)
			if err != nil {
				c.hub.logger.WithError(err).Error("Error marking messages seen")
				continue
			}

			wsMessage.Payload["user_id"] = c.userID
			updatedMessage, err := json.Marshal(wsMessage)
			if err != nil {
				c.hub.logger.WithError(err).Error("Error marshaling updated message")
				continue
			}

			c.hub.SendToConversation(wsMessage.ConversationID, updatedMessage)

		case "reaction_update":
			{
				if _, ok := wsMessage.Payload["reaction"].(map[string]interface{}); ok {
					updatedMessage, err := json.Marshal(wsMessage)
					if err != nil {
						c.hub.logger.WithError(err).Error("Error marshaling reaction update")
						continue
					}
					c.hub.SendToConversation(wsMessage.ConversationID, updatedMessage)
				}
			}
		case "reaction_delete":
			if _, ok := wsMessage.Payload["reaction"].(map[string]interface{}); ok {
				updatedMessage, err := json.Marshal(wsMessage)
				if err != nil {
					c.hub.logger.WithError(err).Error("Error marshaling reaction deletion")
					continue
				}
				c.hub.SendToConversation(wsMessage.ConversationID, updatedMessage)
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		if r := recover(); r != nil {
			c.hub.logger.WithField("recover", r).Error("Recovered in writePump")
		}
		ticker.Stop()
		c.mu.Lock()
		if c.isActive {
			c.isActive = false
			_ = c.conn.Close()
		}
		c.mu.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				// c.hub.logger.WithError(err).Error("Error setting write deadline")
				return
			}
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					// c.hub.logger.WithError(err).Error("Error writing close message")
					return
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				// c.hub.logger.WithError(err).Error("Error writing message")
				return
			}

			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err := w.Write([]byte{'\n'})
				if err != nil {
					c.hub.logger.WithError(err).Error("Error writing newline")
					return
				}
				_, err = w.Write(<-c.send)
				if err != nil {
					c.hub.logger.WithError(err).Error("Error writing newline")
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				c.hub.logger.WithError(err).Error("Error setting write deadline")
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (rt *APIRouter) HandleWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
	defer func() {
		if r := recover(); r != nil {
			c.hub.logger.WithField("recover", r).Error("Recovered in heartbeat")
		}
		ticker.Stop()
	}()

	for {
		c.mu.RLock()
		isActive := c.isActive
		c.mu.RUnlock()

		if !isActive {
			return
		}

		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			c.mu.Lock()
			c.isActive = false
			c.mu.Unlock()
			return
		}
		<-ticker.C
	}
}
