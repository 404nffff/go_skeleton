package handle

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

var clients = make(map[string]*websocket.Conn) // 客户端集合

var msg = make(chan []byte) // 消息通道

func Join(c *gin.Context) {

	username := c.Query("username")

	// if username != "1" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "username is not 1",
	// 	})
	// 	return
	// }

	// 将 HTTP 连接升级为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer ws.Close()

	//加入房间
	clients[username] = ws

	log.Printf("Client connected")

	//打印clients
	fmt.Println("clients:", clients)

	// WebSocket 处理逻辑
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Received message: %s", message)

		// 写入消息通道
		msg <- message
	}
}
