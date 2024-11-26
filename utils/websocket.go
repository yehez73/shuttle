package utils

import (
	"log"
	"shuttle/services"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	activeConnections = make(map[string]*websocket.Conn) // Save active WebSocket connections
	mutex             = &sync.Mutex{}                     // Ensure atomic operations
)

func AddConnection(userID string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	activeConnections[userID] = conn
}

func RemoveConnection(userID string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(activeConnections, userID)
}

func GetConnection(userID string) (*websocket.Conn, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	conn, exists := activeConnections[userID]
	return conn, exists
}

// Handle WebSocket connection
func HandleWebSocketConnection(c *websocket.Conn) {
	userID := c.Params("id")

	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Invalid user ID:", err)
		return
	}

	// Ensure only one connection per user
	if existingConn, exists := GetConnection(userID); exists {
		log.Println("Closing previous connection for user:", userID)
		existingConn.Close()
	}

	AddConnection(userID, c)
	log.Println("User connected:", userID)

	err = services.UpdateUserStatus(ObjectID, "online", time.Time{})
	if err != nil {
		log.Println("Error updating user status:", err)
	}

	// Optional message to client
	err = c.WriteMessage(websocket.TextMessage, []byte("Hello from server!"))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	// Loop to read and write messages
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received message: %s", msg)

		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Println("Error sending response:", err)
			break
		}
	}

	// Disconnect user
	RemoveConnection(userID)
	log.Println("User disconnected:", userID)

	err = services.UpdateUserStatus(ObjectID, "offline", time.Now())
	if err != nil {
		log.Println("Error updating user status:", err)
	}
}