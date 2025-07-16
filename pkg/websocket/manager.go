package websocket

import (
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
}

// ChannelManager manages all channels and clients
type ChannelManager struct {
	clients    map[*Client]bool
	channels   map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	Broadcast  chan Message
	mutex      sync.RWMutex
	shutdown   chan struct{} // Add shutdown channel
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
	for {
		select {
		case client := <-manager.register:
			manager.mutex.Lock()
			manager.clients[client] = true
			manager.clientsMap[client.userId] = client
			manager.mutex.Unlock()
			fmt.Println("New client connected")

		case client := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				delete(manager.clientsMap, client.userId)
				close(client.send)

				// Remove client from all channels
				for channelID := range client.channels {
					if _, ok := manager.channels[channelID]; ok {
						delete(manager.channels[channelID], client)
					}
				}
			}
			manager.mutex.Unlock()
			fmt.Println("Client disconnected")

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

	//fmt.Printf("Client joined channel: %s", channelID)
}

// LeaveChannel removes a client from a specific channel
func (manager *ChannelManager) LeaveChannel(client *Client, channelID string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if _, ok := manager.channels[channelID]; ok {
		delete(manager.channels[channelID], client)
		delete(client.channels, channelID)
		fmt.Printf("Client left channel: %s", channelID)
	}
}

// GetClient Get Client By UserId
func (manager *ChannelManager) GetClient(userId string) *Client {
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
		return true // Message sent successfully
	default:
		// Client's send channel is full or closed
		return false
	}
}

// HandleClient handles WebSocket connections for a client
func (client *Client) HandleClient() {
	defer func() {
		client.channelMgr.unregister <- client
		client.conn.Close()
	}()

	// Handle incoming messages
	go func() {
		for {
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading message: %v", err)
				break
			}

			fmt.Println("Received message:", string(message))

			// Parse the incoming message as JSON
			var action ClientAction
			if err := json.Unmarshal(message, &action); err != nil {
				fmt.Printf("Error parsing message: %v", err)
				continue
			}

			// Handle different actions
			switch action.Action {
			case "subscribe":
				// Join the specified channel
				fmt.Printf("Client subscribing to channel: %s\n", action.ChannelID)
				client.channelMgr.JoinChannel(client, action.ChannelID)

				// Optionally send confirmation back to client
				confirmMsg, _ := json.Marshal(map[string]string{
					"status":     "subscribed",
					"channel_id": action.ChannelID,
				})
				client.conn.WriteMessage(websocket.TextMessage, confirmMsg)

			case "unsubscribe":
				// Leave the specified channel
				fmt.Printf("Client unsubscribing from channel: %s\n", action.ChannelID)
				client.channelMgr.LeaveChannel(client, action.ChannelID)

			case "message":
				// Send message to the specified channel
				client.channelMgr.Broadcast <- Message{
					ChannelID: action.ChannelID,
					Content:   action.Content,
				}

			default:
				fmt.Printf("Unknown action: %s\n", action.Action)
			}
		}
	}()

	// Send messages to this client
	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Printf("Error writing message: %v", err)
			break
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

	client := &Client{
		conn:       ws,
		channels:   make(map[string]bool),
		send:       make(chan []byte, 256),
		channelMgr: manager,
		userId:     userId,
	}

	manager.register <- client

	// Join default channel (you can customize this)
	chatDto := dto.GetChats{}
	chatDto.CauserId = c.Request().Header.Get("causer-id")
	chatDto.CauserType = c.Request().Header.Get("causer-type")
	chats := chatUseCase.GetActiveChats(c.Request().Context(), &chatDto)
	for _, chat := range chats {
		manager.JoinChannel(client, chat.Channel)
	}
	manager.JoinChannel(client, consts.GENERAL_CHAT)

	// Handle client messages
	go client.HandleClient()

	return nil
}
