package controllers

import (
	util "backend_apps/package"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var Connections = make(map[*websocket.Conn]bool)

func SetupWebSocket(app *fiber.App) {
	baseUrl := util.Getenv("MIDDLE_URL", "")
	app.Get(baseUrl+"/ws", websocket.New(HandleWebSocket))
}

func HandleWebSocket(c *websocket.Conn) {
	Connections[c] = true
	defer func() {
		delete(Connections, c)
		c.Close()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}

func Broadcast(message []byte) {
	for conn := range Connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			conn.Close()
			delete(Connections, conn)
		}
	}
}
