package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	web.Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true // Для разработки
	},
}

// @Title ConnectWebSocket
// @Description Establish WebSocket connection
// @Success 200
// @router /ws [get]
func (c *WebSocketController) Get() {
	const writeWait = 10 * time.Second

	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	ws.SetPingHandler(func(appData string) error {
		log.Println("Received ping")
		return ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(writeWait))
	})

	msgChan := make(chan []byte)
	defer close(msgChan)

	// Writer
	go func() {
		for {
			select {
			case message, ok := <-msgChan:
				if !ok {
					return
				}
				if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("Write error:", err)
					return
				}
			}
		}
	}()

	// Reader
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket error:", err)
			}
			break
		}

		log.Printf("Received: %s", message)
		msgChan <- message
	}
}
