package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
)

type Server struct {
	ln      net.Listener
	port    int
	clients []*Client
}

type Client struct {
	id     string
	conn   net.Conn
	server *Server
}

func createServer(port int) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return &Server{ln, port, make([]*Client, 0)}
}

func createClient(conn net.Conn, server *Server) *Client {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	id := fmt.Sprintf("%x", b[:])
	return &Client{id, conn, server}
}

func (s *Server) addConn(conn net.Conn) *Client {
	client := createClient(conn, s)
	s.clients = append(s.clients, client)
	return client
}

func (s *Server) handleConn(c *Client) {
	for {
		msg, _ := bufio.NewReader(c.conn).ReadString('\n')
		fmt.Printf("Message received by %s: %s", c.id, string(msg))
	}
}

func main() {
	server := createServer(8000)
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			panic(err)
		}

		client := server.addConn(conn)
		go server.handleConn(client)
	}
}
