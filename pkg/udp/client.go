package udp

import (
	"fmt"
	"net"
	"tool/global/variable"
)

func Send(message string) {

	ip := variable.ConfigYml.GetString("JobServer.Ip") + ":" + variable.ConfigYml.GetString("JobServer.Port")

	conn, err := net.Dial("udp", ip)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}

	fmt.Println("Message sent:", message)
}
