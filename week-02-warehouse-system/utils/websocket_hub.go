package utils

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var clientsMutex sync.Mutex

func AddClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	log.Println("Klien WebSocket baru terhubung. Total:", len(clients))
}

func RemoveClient(conn *websocket.Conn) {
	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()

	log.Println("Klien WebSocket baru terhubung. Total:", len(clients))
}

func BroadcastNotification(message string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for conn := range clients {
		err := conn.WriteJSON(map[string]interface{}{
			"type":    "NOTIFICATION",
			"message": message,
		})
		if err != nil {
			log.Println("Gagal mengirim ke klien: ", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
