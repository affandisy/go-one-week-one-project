package handlers

import (
	"log"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func WebSocketListen(c *websocket.Conn) {
	utils.AddClient(c)
	defer utils.RemoveClient(c)

	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Koneksi ditutup oleh klien: ", err)
			break
		}

		log.Printf("Pesan diterima: %s (tipe: %d)", msg, messageType)
	}
}
