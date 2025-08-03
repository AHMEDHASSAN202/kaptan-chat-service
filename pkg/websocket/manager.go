package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"kaptan/internal/module/chat/consts"
	"kaptan/internal/module/chat/domain"
	"kaptan/internal/module/chat/dto"
	"kaptan/pkg/utils"
	"net/http"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (customize for production)
	},
}

// Client represents a single websocket connection
type Client struct {
	conn       *websocket.Conn
	channels   map[string]bool
	send       chan []byte
	channelMgr *ChannelManager
	userId     string
	lastSeen   time.Time
	ctx        context.Context
	cancel     context.CancelFunc
}

// ChannelManager manages all channels and clients
type ChannelManager struct {
	clients    map[*Client]bool
	channels   map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Broadcast  chan Message
	mutex      sync.RWMutex
	shutdown   chan struct{}
	clientsMap map[string]*Client
}

// Message represents a message to be sent to a specific channel
type Message struct {
	ChannelID    string `json:"channel_id"`
	Content      string `json:"content"`
	Action       string `json:"action"`
	ExceptClient *Client
}

type ClientAction struct {
	Action    string `json:"action"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content,omitempty"`
	UserId    string `json:"user_id,omitempty"`
}

// Create a new channel manager
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		clients:    make(map[*Client]bool),
		channels:   make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		mutex:      sync.RWMutex{},
		shutdown:   make(chan struct{}),
		clientsMap: make(map[string]*Client),
	}
}

// Run the channel manager
func (manager *ChannelManager) Run() {
	// Start cleanup routine for stale connections
	go manager.cleanupStaleConnections()

	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.clientsMap[client.userId] = client
			client.lastSeen = time.Now()
			manager.mutex.Unlock()
			fmt.Printf("New client connected: %s\n", client.userId)

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				delete(manager.clientsMap, client.userId)
				close(client.send)

				// Cancel client context
				if client.cancel != nil {
					client.cancel()
				}

				// Remove client from all channels
				for channelID := range client.channels {
					if _, ok := manager.channels[channelID]; ok {
						delete(manager.channels[channelID], client)
					}
				}
			}
			manager.mutex.Unlock()
			fmt.Printf("Client disconnected: %s\n", client.userId)

		case message := <-manager.Broadcast:
			manager.mutex.RLock()
			if clients, ok := manager.channels[message.ChannelID]; ok {
				for client := range clients {
					// Skip the sender
					if message.ExceptClient != nil && client == message.ExceptClient {
						continue
					}
					select {
					case client.send <- []byte(utils.JsonEncode(message)):
						client.lastSeen = time.Now() // Update last seen when sending message
					default:
						close(client.send)
						delete(manager.clients, client)
						delete(manager.channels[message.ChannelID], client)
					}
				}
			}
			manager.mutex.RUnlock()

		case <-manager.shutdown:
			// Close all client connections
			manager.mutex.Lock()
			for client := range manager.clients {
				if client.cancel != nil {
					client.cancel()
				}
				close(client.send)
				client.conn.Close()
			}
			manager.clients = make(map[*Client]bool)
			manager.channels = make(map[string]map[*Client]bool)
			manager.mutex.Unlock()
			return // Exit the goroutine
		}
	}
}

// cleanupStaleConnections removes connections that haven't been active
func (manager *ChannelManager) cleanupStaleConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			manager.mutex.Lock()
			now := time.Now()
			var staleClients []*Client

			for client := range manager.clients {
				// Consider connection stale if no activity for 2 minutes
				if now.Sub(client.lastSeen) > 2*time.Minute {
					staleClients = append(staleClients, client)
				}
			}

			// Remove stale clients
			for _, client := range staleClients {
				fmt.Printf("Removing stale client: %s\n", client.userId)
				manager.unregister <- client
			}
			manager.mutex.Unlock()

		case <-manager.shutdown:
			return
		}
	}
}

// JoinChannel adds a client to a specific channel
func (manager *ChannelManager) JoinChannel(client *Client, channelID string) {
	if client == nil {
		return
	}
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	// Initialize the channel if it doesn't exist
	if _, ok := manager.channels[channelID]; !ok {
		manager.channels[channelID] = make(map[*Client]bool)
	}

	// Add client to channel
	manager.channels[channelID][client] = true
	client.channels[channelID] = true
	client.lastSeen = time.Now()

	fmt.Printf("Client %s joined channel: %s\n", client.userId, channelID)
}

// LeaveChannel removes a client from a specific channel
func (manager *ChannelManager) LeaveChannel(client *Client, channelID string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if _, ok := manager.channels[channelID]; ok {
		delete(manager.channels[channelID], client)
		delete(client.channels, channelID)
		fmt.Printf("Client %s left channel: %s\n", client.userId, channelID)
	}
}

// GetClient Get Client By UserId
func (manager *ChannelManager) GetClient(userId string) *Client {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	return manager.clientsMap[userId]
}

func (manager *ChannelManager) Shutdown() {
	close(manager.shutdown)
}

// SendToClient sends a message directly to a specific client by user ID
func (manager *ChannelManager) SendToClient(userId string, message interface{}) bool {
	manager.mutex.RLock()
	client, exists := manager.clientsMap[userId]
	manager.mutex.RUnlock()

	if !exists || client == nil {
		return false // Client not found or offline
	}

	messageBytes := []byte(utils.JsonEncode(message))

	select {
	case client.send <- messageBytes:
		client.lastSeen = time.Now()
		return true // Message sent successfully
	default:
		// Client's send channel is full or closed
		return false
	}
}

// HandleClient handles WebSocket connections for a client with proper timeout handling
func (client *Client) HandleClient() {
	defer func() {
		client.channelMgr.unregister <- client
		client.conn.Close()
	}()

	// Set up ping/pong handlers
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		client.lastSeen = time.Now()
		return nil
	})

	// Start ping routine
	go client.pingRoutine()

	// Handle incoming messages
	go client.readMessages()

	// Send messages to this client
	client.writeMessages()
}

// readMessages handles incoming messages from the client
func (client *Client) readMessages() {
	defer func() {
		if client.cancel != nil {
			client.cancel()
		}
	}()

	for {
		select {
		case <-client.ctx.Done():
			return
		default:
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket error for client %s: %v\n", client.userId, err)
				}
				return
			}

			client.lastSeen = time.Now()
			fmt.Printf("Received message from %s: %s\n", client.userId, string(message))

			// Parse the incoming message as JSON
			var action ClientAction
			if err := json.Unmarshal(message, &action); err != nil {
				fmt.Printf("Error parsing message from %s: %v\n", client.userId, err)
				continue
			}

			// Handle different actions
			switch action.Action {
			case "subscribe":
				fmt.Printf("Client %s subscribing to channel: %s\n", client.userId, action.ChannelID)
				client.channelMgr.JoinChannel(client, action.ChannelID)

				// Send confirmation back to client
				confirmMsg := map[string]interface{}{
					"type":       "subscription_confirmed",
					"status":     "subscribed",
					"channel_id": action.ChannelID,
					"timestamp":  time.Now().Unix(),
				}

				if confirmBytes, err := json.Marshal(confirmMsg); err == nil {
					select {
					case client.send <- confirmBytes:
					default:
						return
					}
				}

			case "unsubscribe":
				fmt.Printf("Client %s unsubscribing from channel: %s\n", client.userId, action.ChannelID)
				client.channelMgr.LeaveChannel(client, action.ChannelID)

			case "message":
				client.channelMgr.Broadcast <- Message{
					ChannelID: action.ChannelID,
					Content:   action.Content,
				}

			case "ping":
				// Handle explicit ping from client
				pongMsg := map[string]interface{}{
					"type":      "pong",
					"timestamp": time.Now().Unix(),
				}
				if pongBytes, err := json.Marshal(pongMsg); err == nil {
					select {
					case client.send <- pongBytes:
					default:
						return
					}
				}

			default:
				fmt.Printf("Unknown action from %s: %s\n", client.userId, action.Action)
			}
		}
	}
}

// writeMessages handles outgoing messages to the client
func (client *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Printf("Error writing message to client %s: %v\n", client.userId, err)
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("Error sending ping to client %s: %v\n", client.userId, err)
				return
			}

		case <-client.ctx.Done():
			return
		}
	}
}

// pingRoutine sends periodic ping messages to keep connection alive
func (client *Client) pingRoutine() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send application-level ping
			pingMsg := map[string]interface{}{
				"type":      "ping",
				"timestamp": time.Now().Unix(),
			}

			if pingBytes, err := json.Marshal(pingMsg); err == nil {
				select {
				case client.send <- pingBytes:
				case <-client.ctx.Done():
					return
				default:
					// Channel is full, connection might be dead
					return
				}
			}

		case <-client.ctx.Done():
			return
		}
	}
}

// WebSocket handler function
func handleWebSocket(c echo.Context, manager *ChannelManager, chatUseCase domain.ChatRepository) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	userId := utils.GetClientUserId(c.Request().Header.Get("causer-type"), c.Request().Header.Get("causer-id"))

	// Create context with cancel for the client
	ctx, cancel := context.WithCancel(c.Request().Context())

	client := &Client{
		conn:       ws,
		channels:   make(map[string]bool),
		send:       make(chan []byte, 256),
		channelMgr: manager,
		userId:     userId,
		lastSeen:   time.Now(),
		ctx:        ctx,
		cancel:     cancel,
	}

	manager.register <- client

	// Join default channels
	chatDto := dto.GetChats{}
	chatDto.CauserId = c.Request().Header.Get("causer-id")
	chatDto.CauserType = c.Request().Header.Get("causer-type")
	chats := chatUseCase.GetActiveChats(c.Request().Context(), &chatDto)
	for _, chat := range chats {
		manager.JoinChannel(client, chat.Channel)
	}
	manager.JoinChannel(client, consts.GENERAL_CHAT)

	// Handle client messages (this will block until client disconnects)
	client.HandleClient()

	return nil
}
