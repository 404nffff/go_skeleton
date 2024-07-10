package tcp

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// TCPServer 结构体表示TCP服务器
type TCPServer struct {
	address string
}

// NewTCPServer 创建一个新的TCPServer
func NewTCPServer(address string) *TCPServer {
	return &TCPServer{address: address}
}

// Start 启动TCP服务器
func (s *TCPServer) Start() {
	// 监听指定的端口
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on", s.address)

	for {
		// 接受来自客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go s.handleRequest(conn)
	}
}

// handleRequest 处理客户端请求
func (s *TCPServer) handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	fmt.Printf("Message received: %s", message)

	// 回复客户端
	_, err = conn.Write([]byte("Message received\n"))
	if err != nil {
		fmt.Println("Error writing:", err.Error())
	}

	//如果message是stop
	if message == "stop\n" {
		//打印
		fmt.Println("Stop server")
		//(event_manage.CreateEventManageFactory()).FuzzyCall(variable.EventDestroyPrefix)
	}
}
