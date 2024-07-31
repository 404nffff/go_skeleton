package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8081/join"

	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Printf("Connected to %s", serverAddr)

	// 创建一个通道来接收操作系统的中断信号
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 创建一个通道用于从标准输入读取消息
	messageChan := make(chan string)

	// 在一个新的 goroutine 中从标准输入读取消息
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter message: ")
			message, _ := reader.ReadString('\n')
			messageChan <- message
		}
	}()

	// 在另一个 goroutine 中接收服务器消息
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	// 主循环
	for {
		select {
		case <-interrupt:
			log.Println("Interrupt received, closing connection")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			return
		case message := <-messageChan:
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
