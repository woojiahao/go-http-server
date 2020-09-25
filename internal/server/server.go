package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/woojiahao/go-http-server/internal/client"
	"net"
	"strings"
)

type Server struct {
	Ln      net.Listener
	port    int
	clients []*client.Client
	done    chan bool
}

func Create(port int) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	return &Server{
		ln,
		port,
		make([]*client.Client, 0),
		make(chan bool, 1),
	}
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

	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(scanLinesWithCR)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println(message)
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

func (s *Server) Start() {
	fmt.Printf("Creating server on port %d\n", s.port)
	fmt.Printf("http://127.0.0.1:%d\n", s.port)

	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("Connection failed: %v", err)
				return
			}
		}

		client := s.AddConn(conn)
		fmt.Printf("New connection made to client %s\n", client.ID)
		go s.HandleConn(client)
	}
}

func (s *Server) Stop() {
	fmt.Println("Stopping server")
	s.done <- true
	s.Ln.Close()
}

func scanLinesWithCR(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func processMessage(message string) (string, error) {
	lines := strings.Split(message, "\n")
	parts := strings.Split(lines[0], " ")
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
