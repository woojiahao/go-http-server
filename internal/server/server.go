package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/woojiahao/go-http-server/internal/client"
	"net"
	"os"
	"strings"
)

type Server struct {
	Ln      net.Listener
	port    int
	clients []*client.Client
	done    chan bool
	path    string
}

func Create(port int) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &Server{
		ln,
		port,
		make([]*client.Client, 0),
		make(chan bool, 1),
		path,
	}
}

func (s *Server) AddConn(conn net.Conn) *client.Client {
	client := client.Create(conn)
	s.clients = append(s.clients, client)
	return client
}

func (s *Server) HandleConn(c *client.Client) {
	defer func() {
		fmt.Printf("Connection %s dropped\n", c.ID)
		_ = c.Conn.Close()
	}()

	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(scanLinesWithCR)

	// First line is the start line
	if scanner.Scan() {
		startLine := scanner.Text()
		method, resource, httpVersion, err := parseStartLine(s.path, startLine)
		if err != nil {
			fmt.Printf("Invalid request: %s\n", err.Error())
			c.Conn.Write([]byte(err.Error()))
		}
		fmt.Printf("\t-- Method:\t%s\n", string(method))
		fmt.Printf("\t-- Resource:\t%s\n", resource)
		fmt.Printf("\t-- HTTP Ver:\t%s\n", httpVersion)
	}

	// Second line onwards is the headers
	// TODO Check why headers are returned on new lines rather than all together and then with a \r\n at the end like the
	// document suggest
	headers := make(map[string]string)
	for scanner.Scan() {
		h := strings.TrimSpace(scanner.Text())
		if h == "" {
			break
		}
		key, value, err := parseHeader(h)
		if err != nil {
			fmt.Printf("Invalid header: %s\n", h)
			c.Conn.Write([]byte(err.Error()))
		}
		headers[key] = value
	}
	fmt.Println(headers)
}

func (s *Server) Start() {
	fmt.Println("Server start")
	fmt.Printf("\tCreating server on port %d\n", s.port)
	fmt.Printf("\thttp://127.0.0.1:%d\n", s.port)
	path, _ := os.Getwd()
	fmt.Printf("\tCurrent directory: %s\n", path)

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
	// TODO Figure out how this works
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
