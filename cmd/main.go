package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strings"
)

type Server struct {
	ln      net.Listener
	port    int
	clients []*Client
	backlog int
}

type Client struct {
	id     string
	conn   net.Conn
	server *Server
}

func createServer(port int, backlog int) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return &Server{ln, port, make([]*Client, 0), backlog}
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
	defer func() {
		_ = c.conn.Close()
	}()

	for {
		content := make([]string, 0)
		for {
			// Keep reading the input from the client and adding it to the content
			msg, _, err := bufio.NewReader(c.conn).ReadLine()
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", c.id)
				return
			}
			if string(msg) == "" {
				break
			}
			content = append(content, string(msg))
		}
		message := strings.Join(content, "\n")
		fmt.Printf("Message received by %s: %s\n", c.id, message)
		c.conn.Write([]byte(fmt.Sprintf("-> You sent %s\n", message)))
	}
}

func main() {
	server := createServer(8000, 1)

	for {
		conn, err := server.ln.Accept()
		if err != nil {
			panic(err)
		}

		client := server.addConn(conn)
		fmt.Printf("New connection made to client %s\n", client.id)
		go server.handleConn(client)
	}
}
