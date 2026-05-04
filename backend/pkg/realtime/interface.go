package realtime

import "github.com/gofiber/contrib/websocket"

type IHub interface {
	Run()
	Subscribe(topic string, conn *websocket.Conn)
	Unsubscribe(topic string, conn *websocket.Conn)
	Broadcast(topic string, payload interface{})
}
