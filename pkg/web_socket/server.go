package web_socket

import (
	"fmt"
	"sync"
	"time"
)

var (
	h    *Hub
	once sync.Once
)

func NewHub() *Hub {
	once.Do(func() {
		h = &Hub{
			broadcast:  make(chan []byte),
			clients:    make(map[*Client]int),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
	})
	return h
}

// 运行
func (h *Hub) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register: // 注册
			h.clients[client] = 1
		case client := <-h.unregister: // 注销

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		case <-ticker.C:
			// 可以考虑减少日志频率或只在调试模式下记录
			//log.Println("No message in the last second")

			fmt.Println(h.clients)
		}
	}
}
