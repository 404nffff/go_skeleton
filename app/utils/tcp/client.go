package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
)

// TCPClient 结构体表示TCP客户端
type TCPClient struct {
	address string
}

// NewTCPClient 创建一个新的TCPClient
func NewTCPClient(address string) *TCPClient {
	return &TCPClient{address: address}
}

// Connect 连接到服务器并发送消息
func (c *TCPClient) Connect(message string) {
	// 连接到服务器
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// 发送消息到服务器
	_, err = conn.Write(bytes.NewBufferString(message + "\n").Bytes())
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
		return
	}

	// 读取服务器的响应（如果有）
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err.Error())
		return
	}

	fmt.Printf("Response from server: %s\n", response)
}
