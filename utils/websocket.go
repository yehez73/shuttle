package utils

import (
	"log"
	"sync"
	"shuttle/services"

	"github.com/gofiber/contrib/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	activeConnections = make(map[string]*websocket.Conn) // menyimpan koneksi aktif
	mutex             = &sync.Mutex{}                     // memastikan akses ke activeConnections aman
)

// Handle koneksi WebSocket dan update status online/offline
func HandleWebSocketConnection(c *websocket.Conn) {
	userID := c.Params("id")

	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Invalid user ID:", err)
		return
	}

	// Lock untuk memastikan data tidak konflik
	mutex.Lock()
	defer mutex.Unlock()

	// Menangani koneksi baru untuk user yang terhubung
	if existingConn, exists := activeConnections[userID]; exists {
		// Jika ada koneksi lama, tutup koneksi sebelumnya
		log.Println("Closing previous connection for user:", userID)
		existingConn.Close()
	}

	// Simpan koneksi baru
	activeConnections[userID] = c
	log.Println("User connected:", userID)

	// Update status online di database
	err = services.UpdateUserStatus(ObjectID, "online")
	if err != nil {
		log.Println("Error updating user status:", err)
	}

	// Kirim pesan selamat datang
	err = c.WriteMessage(websocket.TextMessage, []byte("Hello from server!"))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	// Terima pesan dari client dan kirim kembali (echo)
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received message: %s", msg)

		// Kirim kembali pesan ke client
		err = c.WriteMessage(mt, msg)
		if err != nil {
			log.Println("Error sending response:", err)
			break
		}
	}

	// Menghapus koneksi ketika WebSocket ditutup
	delete(activeConnections, userID)
	log.Println("User disconnected:", userID)

	// Update status offline di database
	err = services.UpdateUserStatus(ObjectID, "offline")
	if err != nil {
		log.Println("Error updating user status:", err)
	}
}