package web_socket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// 客户端合集
	clients map[*Client]int

	// Inbound messages from the clients.
	broadcast chan []byte

	//注册
	register chan *Client

	//注销
	unregister chan *Client
}

// 客户端
type Client struct {
	// 客户端连接
	conn *websocket.Conn

	// 发送通道
	send chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
	ReadBufferSize:  1024, // 读取缓冲区大小
	WriteBufferSize: 1024, // 写入缓冲区大小
}
