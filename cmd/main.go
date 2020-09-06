package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
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
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == io.EOF {
			fmt.Printf("Client %s disconnected", c.id)
			return
		}
		fmt.Printf("Message received by %s: %s", c.id, string(msg))
		c.conn.Write([]byte(fmt.Sprintf("-> You sent %s", string(msg))))
	}
}

// Monitors the backlog of the server. Every time a new connection comes in and the backlog is not full, the connection
// will be picked off from the channel and processed. A new client is created and added to the server to manage.
func monitorBacklog(server *Server, backlog <-chan net.Conn) {
	for {
		conn := <-backlog
		client := server.addConn(conn)
		fmt.Printf("New connection made to client %s\n", client.id)
		go server.handleConn(client)
	}
}

func main() {
	server := createServer(8000, 1)
	// Create a backlog to house the pending connection requests that can be held at once
	backlog := make(chan net.Conn, server.backlog)

	go monitorBacklog(server, backlog)

	for {
		conn, err := server.ln.Accept()
		if err != nil {
			panic(err)
		}

		select {
		case backlog <- conn:
		default:
			// If the channel is already full
			fmt.Println("Connection rejected")
			conn.Write([]byte("Queue is full!"))
			conn.Close()
		}
	}
}
