package utils

import (
	"encoding/json"
	"sync"
	"time"

	"shuttle/logger"
	"shuttle/repositories"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketServiceInterface interface {
	HandleWebSocketConnection(c *websocket.Conn)
}

type WebSocketService struct {
	userRepository repositories.UserRepositoryInterface
	authRepository repositories.AuthRepositoryInterface
}

func NewWebSocketService(userRepository repositories.UserRepositoryInterface, authRepository repositories.AuthRepositoryInterface) WebSocketServiceInterface {
	return &WebSocketService{
		userRepository: userRepository,
		authRepository: authRepository,
	}
}

var (
	activeConnections = make(map[string]*websocket.Conn) // Save active WebSocket connections
	mutex             = &sync.Mutex{}                    // Ensure atomic operations
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
func (s *WebSocketService) HandleWebSocketConnection(c *websocket.Conn) {
	UUID := c.Params("id")

	_, err := s.userRepository.FetchSpecificUser(UUID)
	if err != nil {
		logger.LogError(err, "Websocket Error Getting User", nil)
		return
	}

	// Ensure only one connection per user
	if existingConn, exists := GetConnection(UUID); exists {
		logger.LogInfo("Websocket Connection Already Exists, Closing Existing Connection", map[string]interface{}{"ID": UUID})
		existingConn.Close()
	}

	AddConnection(UUID, c)
	logger.LogInfo("Websocket Connection Established", map[string]interface{}{"ID": UUID})

	err = s.authRepository.UpdateUserStatus(UUID, "online", time.Time{})
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

		logger.LogInfo("Websocket Message Parsed", map[string]interface{}{"UUID": UUID, "longitude": data.Longitude, "latitude": data.Latitude})

		response := struct {
			Code    int    `json:"code"`
			Status  string `json:"status"`
			Message string `json:"message"`
		}{
			Code:    200,
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
	RemoveConnection(UUID)
	logger.LogInfo("Websocket Connection Closed", map[string]interface{}{"ID": UUID})

	err = s.authRepository.UpdateUserStatus(UUID, "offline", time.Now())
	if err != nil {
		logger.LogError(err, "Websocket Error Updating User Status", nil)
	}
}
