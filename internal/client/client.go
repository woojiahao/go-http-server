package client

import (
	"fmt"
	"math/rand"
	"net"
)

type Client struct {
	ID   string
	Conn net.Conn
}

func Create(conn net.Conn) *Client {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	id := fmt.Sprintf("%x", b[:])
	return &Client{id, conn}
}
