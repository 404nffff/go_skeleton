package web_socket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader  = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
	ReadBufferSize:  1024, // 读取缓冲区大小
	WriteBufferSize: 1024, // 写入缓冲区大小
}
