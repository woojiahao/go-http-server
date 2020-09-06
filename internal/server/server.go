package server

import (
	"fmt"
	"net"
)

func CreateSocket(port int) net.Conn {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	conn, _ := ln.Accept()
	return conn
}
