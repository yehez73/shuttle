package utils

import (
	"encoding/json"
	"sync"
	"time"

	"shuttle/logger"
	"shuttle/repositories"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type WebSocketServiceInterface interface {
	HandleWebSocketConnection(c *websocket.Conn)
}

type WebSocketService struct {
	userRepository repositories.UserRepositoryInterface
	authRepository repositories.AuthRepositoryInterface
	shuttleRepository repositories.ShuttleRepositoryInterface
}

func NewWebSocketService(userRepository repositories.UserRepositoryInterface, authRepository repositories.AuthRepositoryInterface, shuttleRepository repositories.ShuttleRepositoryInterface) WebSocketServiceInterface {
	return &WebSocketService{
		userRepository: userRepository,
		authRepository: authRepository,
		shuttleRepository: shuttleRepository,
	}
}

type WebSocketResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	activeConnections = make(map[string]*websocket.Conn)
	mutex             = &sync.Mutex{}

	connGroups = make(map[string]map[string]*websocket.Conn)
	groupMutex = &sync.Mutex{}
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

func AddToGroup(shuttleUUID, userUUID string, conn *websocket.Conn) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if _, exists := connGroups[shuttleUUID]; !exists {
		connGroups[shuttleUUID] = make(map[string]*websocket.Conn)
	}
	connGroups[shuttleUUID][userUUID] = conn
}

func RemoveFromGroup(shuttleUUID, userUUID string) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if group, exists := connGroups[shuttleUUID]; exists {
		delete(group, userUUID)
		if len(group) == 0 {
			delete(connGroups, shuttleUUID)
		}
	}
}

func BroadcastToGroup(shuttleUUID string, message []byte) {
	groupMutex.Lock()
	defer groupMutex.Unlock()

	if group, exists := connGroups[shuttleUUID]; exists {
		for _, conn := range group {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.LogError(err, "WebSocket Broadcast Error", nil)
			}
		}
	}
}

// Handle WebSocket connection
func (s *WebSocketService) HandleWebSocketConnection(c *websocket.Conn) {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		response := WebSocketResponse{
			Code:    401,
			Status:  "Unauthorized",
			Message: "Unauthorized access",
		}
		responseMsg, _ := json.Marshal(response)
		c.WriteMessage(websocket.TextMessage, responseMsg)
		c.Close()
		return
	}
	shuttleUUID := c.Params("id")

	userUUIDParsed, err := uuid.Parse(userUUID)
	if err != nil {
		logger.LogError(err, "Invalid UUID format", nil)
		c.Close()
		return
	}

	err = s.userRepository.UpdateUserStatus(userUUIDParsed, "online", time.Time{})
	if err != nil {
		logger.LogError(err, "WebSocket Error Updating User Status", nil)
	}

	defer func() {
		if shuttleUUID != "" {
			RemoveFromGroup(shuttleUUID, userUUID)
		} else {
			RemoveConnection(userUUID)
		}

		if err := s.userRepository.UpdateUserStatus(userUUIDParsed, "offline", time.Now()); err != nil {
			logger.LogError(err, "WebSocket Error Updating User Status", nil)
		}

		logger.LogInfo("WebSocket connection closed", map[string]interface{}{
			"UserUUID":    userUUID,
			"ShuttleUUID": shuttleUUID,
		})
	}()

	if shuttleUUID != "" {
		shuttleUUIDParsed, err := uuid.Parse(shuttleUUID)
		if err != nil {
			logger.LogError(err, "Invalid UUID format", nil)
			c.Close()
			return
		}

		exist, err := s.shuttleRepository.CheckIfExistInShuttle(userUUIDParsed, shuttleUUIDParsed)
		if err != nil {
			c.Close()
			return
		}

		if !exist {
			response := WebSocketResponse{
				Code:    404,
				Status:  "Not Found",
				Message: "User not found in shuttle",	
			}
			responseMsg, _ := json.Marshal(response)
			c.WriteMessage(websocket.TextMessage, responseMsg)
			
			logger.LogError(err, "User not found in shuttle", map[string]interface{}{
				"ShuttleUUID": shuttleUUID,
				"UserUUID":    userUUID,
			})

			c.Close()
			return
		}

		AddToGroup(shuttleUUID, userUUID, c)
		logger.LogInfo("WebSocket Connection Added to Group", map[string]interface{}{
			"ShuttleUUID": shuttleUUID,
			"UserUUID":    userUUID,
		})
		c.WriteMessage(websocket.TextMessage, []byte("Connected to group"))
	} else {
		AddConnection(userUUID, c)
		logger.LogInfo("WebSocket connection established for user", map[string]interface{}{
			"UserUUID": userUUID,
		})
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		if shuttleUUID != "" {
			var data struct {
				Longitude float64 `json:"longitude"`
				Latitude  float64 `json:"latitude"`
			}

			if err := json.Unmarshal(msg, &data); err != nil || data.Longitude == 0 || data.Latitude == 0 {
				response := WebSocketResponse{
					Code:    400,
					Status:  "Bad Request",
					Message: "Invalid request format",
				}
				responseMsg, _ := json.Marshal(response)
				c.WriteMessage(websocket.TextMessage, responseMsg)
				continue
			}

			logger.LogInfo("Broadcasting Message", map[string]interface{}{
				"ShuttleUUID": shuttleUUID,
				"UserUUID":    userUUID,
				"Longitude":   data.Longitude,
				"Latitude":    data.Latitude,
			})
			BroadcastToGroup(shuttleUUID, msg)

			response := WebSocketResponse{
				Code:    200,
				Status:  "OK",
				Message: "Message broadcasted",
			}
			responseMsg, _ := json.Marshal(response)
			c.WriteMessage(websocket.TextMessage, responseMsg)
		}
	}
}