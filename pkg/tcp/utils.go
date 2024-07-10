package tcp

import "net"

func IsPortInUse(port string) bool {
	conn, err := net.Listen("tcp", port)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}
