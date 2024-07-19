package main

import (
	"fmt"
	"net"
	"os"
	"tool/bootstrap"
	"tool/app/global/variable"
)

func main() {
	
	// 初始化全局变量
	bootstrap.Initialize()

	port := variable.ConfigYml.GetInt("JobServer.Port")
	ip := variable.ConfigYml.GetString("JobServer.Ip")

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server is listening on port", port)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		fmt.Printf("Received message from %v: %s\n", addr, string(buf[:n]))
	}
}
