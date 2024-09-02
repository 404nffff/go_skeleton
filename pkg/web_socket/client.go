package web_socket

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 加入
func Join(w http.ResponseWriter, r *http.Request) *Client {

	// Upgrade the HTTP server connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return nil
	}

	conn.EnableWriteCompression(true)

	// 初始超时时间
	//timeout := 100 * time.Millisecond
	maxTimeout := 5 * time.Second

	conn.SetReadDeadline(time.Now().Add(maxTimeout))

	conn.SetCompressionLevel(9)

	conn.SetPingHandler(func(appData string) error {

		fmt.Println("ping: ", appData)

		conn.SetReadDeadline(time.Now().Add(maxTimeout))

		return nil
	})

	conn.SetPongHandler(func(appData string) error {
		fmt.Println("pong: ", appData)
		return nil
	})

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Println("close: ", code, text)
		return nil
	})

	// Register the client
	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	fmt.Println("client: ", client)

	h.register <- client

	fmt.Println("client: ", client)

	go client.ReadChannel()

	go client.SendChannel()

	return client
}

// 接收消息通道
func (c *Client) ReadChannel() {

	defer func() {
		h.unregister <- c
	}()

	for {

		messageType, reader, err := c.conn.NextReader()

		fmt.Println("type1: ", messageType, reader, err)

		if err != nil {
			h.unregister <- c
			break
		}

		message, err := io.ReadAll(reader)

		if err != nil {
			h.unregister <- c
			fmt.Println("read error: ", err)
			return
		}

		// _, message, err := c.conn.ReadMessage()

		fmt.Println("message: ", string(message), err)

		h.broadcast <- message
	}
}

// 发送消息
func (c *Client) SendChannel() {

	for {
		select {
		case msg, ok := <-c.send:

			if !ok {
				fmt.Println("send channel closed")
				return
			}

			fmt.Println("send: ", string(msg))

			writer, err := c.conn.NextWriter(websocket.TextMessage)	

			if err != nil {

				h.unregister <- c

				fmt.Println("next writer error: ", err)
				return
			}

			writer.Write(msg)

			writer.Close()
		}
	}
}
