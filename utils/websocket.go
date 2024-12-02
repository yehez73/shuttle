package utils

import (
	"encoding/json"
	"shuttle/logger"
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

func AddConnection(ID string, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	activeConnections[ID] = conn
}

func RemoveConnection(ID string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(activeConnections, ID)
}

func GetConnection(ID string) (*websocket.Conn, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	conn, exists := activeConnections[ID]
	return conn, exists
}

// Handle WebSocket connection
func HandleWebSocketConnection(c *websocket.Conn) {
	ID := c.Params("id")

	_, err := services.GetSpecUser(ID)
	if err != nil {
		logger.LogError(err, "Websocket Error Getting User", map[string]interface{}{"ID": ID})
		return
	}

	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		logger.LogError(err, "Websocket Error Parsing ObjectID", nil)
		return
	}

	// Ensure only one connection per user
	if existingConn, exists := GetConnection(ID); exists {
		logger.LogInfo("Websocket Connection Already Exists, Closing Existing Connection", map[string]interface{}{"ID": ID})
		existingConn.Close()
	}

	AddConnection(ID, c)
	logger.LogInfo("Websocket Connection Established", map[string]interface{}{"ID": ID})

	err = services.UpdateUserStatus(ObjectID, "online", time.Time{})
	if err != nil {
		logger.LogError(err, "Websocket Error Updating User Status", nil)
	}

	err = c.WriteMessage(websocket.TextMessage, []byte("Connected to websocket"))
	if err != nil {
		logger.LogError(err, "Websocket Error Writing Message", nil)
		return
	}

	// Loop to read and write messages
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			logger.LogError(err, "Websocket Error Reading Message", nil)
			break
		}
	
		var data struct {
			Longitude float64 `json:"longitude"`
			Latitude  float64 `json:"latitude"`
		}
		
		if err := json.Unmarshal(msg, &data); err != nil {
			logger.LogError(err, "Websocket Message Received Is Not A Location", nil)
			break
		}
	
		logger.LogInfo("Websocket Message Parsed", map[string]interface{}{"ID": ID, "longitude": data.Longitude, "latitude": data.Latitude})
	
		response := struct {
			Status  string  `json:"status"`
			Message string  `json:"message"`
		}{
			Status:  "OK",
			Message: "Data received successfully",
		}
	
		responseMsg, err := json.Marshal(response)
		if err != nil {
			logger.LogError(err, "Error marshaling response message", nil)
			break
		}
	
		err = c.WriteMessage(mt, responseMsg)
		if err != nil {
			logger.LogError(err, "Websocket Error Writing Message", nil)
			break
		}
	}	

	// Disconnect user
	RemoveConnection(ID)
	logger.LogInfo("Websocket Connection Closed", map[string]interface{}{"ID": ID})

	err = services.UpdateUserStatus(ObjectID, "offline", time.Now())
	if err != nil {
		logger.LogError(err, "Websocket Error Updating User Status", nil)
	}
}