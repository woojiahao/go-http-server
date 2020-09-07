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
		c.Conn.Write([]byte(fmt.Sprintf("-- You sent %s\n", message)))
		res, err := processMessage(message)
		if err != nil {
			c.Conn.Write([]byte(fmt.Sprintf("!- ERROR %s\n", err)))
		} else {
			c.Conn.Write([]byte(fmt.Sprintf(">- %s\n", res)))
		}
	}
}

func processMessage(message string) (string, error) {
	parts := strings.Split(message, " ")
	protocol := Keyword(parts[0])

	switch protocol {
	case GET:
		word, err := handleGET(parts[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("ANSWER %s", word), nil

	case SET:
		handleSET(parts[1], strings.Join(parts[2:], "\n"))
		return formatData(), nil

	case CLEAR:
		handleCLEAR()
		return formatData(), nil

	case ALL:
		output := handleALL()
		return strings.Join(output, "\n"), nil

	default:
		return "", fmt.Errorf("invalid protocol (%s) used", protocol)
	}

	return "", nil
}
