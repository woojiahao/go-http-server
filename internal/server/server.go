package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/woojiahao/go-http-server/internal/client"
	"net"
	"os"
	"path/filepath"
	"regexp"
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
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Message received by %s: %s\n", c.ID, message)
		method, resource, httpVersion, err := extractHTTPRequest(s.path, message)
		if err != nil {
			fmt.Printf("Invalid request: %s\n", err.Error())
			c.Conn.Write([]byte(err.Error()))
		}
		fmt.Printf("\t-- Method:\t%s\n", string(method))
		fmt.Printf("\t-- Resource:\t%s\n", resource)
		fmt.Printf("\t-- HTTP Ver:\t%s\n", httpVersion)
		break
	}
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

func extractHTTPRequest(path, request string) (method Method, resource string, httpVersion string, err error) {
	lines := strings.Split(request, "\n")
	parts := strings.Split(lines[0], " ")

	if len(parts) != 3 {
		return Method(""), "", "", fmt.Errorf("HTTP request must include [method] [resource] [http-version]\\r\\n")
	}

	method, resource, httpVersion = Method(parts[0]), parts[1], parts[2]
	if !method.isValid() {
		err = fmt.Errorf("Invalid method. Methods available: %v", methods)
		return
	}

	// TODO Allow users to customise the folder to serve
	fmt.Println(filepath.Join(path, resource))
	if _, e := os.Stat(filepath.Join(path, resource)); os.IsNotExist(e) {
		err = fmt.Errorf("Invalid resource.")
		return
	}

	if match, _ := regexp.MatchString("^HTTP/(0.9|1.0|1.1|2.0)$", httpVersion); !match {
		err = fmt.Errorf("Invalid HTTP version. Available versions: [0.9, 1.0, 1.1, 2.0]")
		return
	}

	return
}
