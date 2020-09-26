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

func (s *Server) HandleConn(c *client.Client) {
	defer func() {
		fmt.Printf("Connection %s closed\n", c.ID)
		_ = c.Conn.Close()
	}()

	// Create a scanner that only stops reading when CRLF is read
	scanner := bufio.NewScanner(c.Conn)
	scanner.Split(scanLinesWithCR)

	// Read the request from the client
	request, err := s.readRequest(scanner)
	if err != nil {
		fmt.Printf("Invalid request: %s\n", err.Error())
		c.Conn.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("%s request for %s on %s\n", string(request.method), request.resource, request.httpVersion)
	fmt.Printf("Headers: %v\n", request.headers)

}

func (s *Server) readRequest(scanner *bufio.Scanner) (request Request, err error) {
	headers := make(map[string]string)
	var method Method
	var resource string
	var httpVersion string

	// First line is the start line
	if scanner.Scan() {
		startLine := scanner.Text()
		method, resource, httpVersion, err = parseStartLine(s.path, startLine)
		if err != nil {
			return
		}
	}

	// Second line onwards is the headers
	// TODO Check why headers are returned on new lines rather than all together and then with a \r\n at the end like the
	// document suggest
	for scanner.Scan() {
		h := strings.TrimSpace(scanner.Text())
		if h == "" {
			break
		}
		key, value, e := parseHeader(h)
		if e != nil {
			err = e
			return
		}
		headers[key] = value
	}

	request = Request{method, resource, httpVersion, headers}

	return
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
