package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8080/ws"

	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// 在连接建立后等待一小段时间再发送初始消息
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		log.Println("write initial message:", err)
		return
	}

	go func() {

		//接收消息
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}

	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 发送 ping 消息
			err := c.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Println("write ping:", err)
				return
			}

			err = c.WriteMessage(websocket.TextMessage, []byte("hello"))

			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
