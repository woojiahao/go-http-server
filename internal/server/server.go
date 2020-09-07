package server

import (
	"bufio"
	"fmt"
	"github.com/woojiahao/go-http-server/internal/client"
	"io"
	"net"
	"strings"
)

type Server struct {
	Ln      net.Listener
	port    int
	clients []*client.Client
}

func Create(port int) *Server {
	fmt.Printf("Creating server on port %d\n", port)
	fmt.Printf("http://127.0.0.1:%d\n", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return &Server{ln, port, make([]*client.Client, 0)}
}

func (s *Server) AddConn(conn net.Conn) *client.Client {
	client := client.Create(conn)
	s.clients = append(s.clients, client)
	return client
}

func (s *Server) HandleConn(c *client.Client) {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		content := make([]string, 0)
		for {
			// Keep reading the input from the client and adding it to the content
			msg, _, err := bufio.NewReader(c.Conn).ReadLine()
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", c.ID)
				return
			}
			if string(msg) == "" {
				break
			}
			content = append(content, string(msg))
		}
		message := strings.Join(content, "\n")
		fmt.Printf("Message received by %s: %s\n", c.ID, message)
		c.Conn.Write([]byte(fmt.Sprintf("-> You sent %s\n", message)))

		// Get the protocol
		protocol := Keyword(strings.Split(message, " ")[0])
		switch protocol {
		case GET:
			fmt.Printf("GET protocol received")
		}
	}
}
