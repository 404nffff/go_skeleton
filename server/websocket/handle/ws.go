package handle

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func HandleMsg() {
	for {
		// 从消息通道中读取消息
		message := <-msg

		fmt.Println("message:", string(message))

		// 给所有客户端发送消息
		for username, client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, username)
			}

			// 判断是否是关闭消息
			if string(message) == "close" {
				client.Close()
				delete(clients, username)
			}
		}
	}
}
