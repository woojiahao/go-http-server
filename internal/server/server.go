package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Server struct {
	ln   net.Listener
	done chan bool
	port int
	path string
}

func Create(port int, path string) *Server {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	return &Server{ln, make(chan bool, 1), port, path}
}

func (s *Server) Start() {
	fmt.Printf("Creating server on port %d\n", s.port)
	fmt.Printf("http://127.0.0.1:%d\n", s.port)

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("Connection failed: %v", err)
				return
			}
		}

		go s.HandleConn(conn)
	}
}

func (s *Server) Stop() {
	fmt.Println("Stopping server")
	s.done <- true
	s.ln.Close()
}

func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	// Create a scanner that only stops reading when CRLF is read
	scanner := bufio.NewScanner(conn)
	scanner.Split(scanLinesWithCR)

	// Read the request from the client
	request, err := readRequest(scanner, s.path)
	if err != nil {
		fmt.Printf("Invalid request: %s\n", err.Error())
		conn.Write([]byte(err.Error()))
		return
	}

	response := generateResponse(request, s.path)
	conn.Write([]byte(response.Serialize()))

	fmt.Printf("%s :: %s :: %d\n", string(request.method), request.resource, response.statusCode.code)
}

func generateResponse(request Request, path string) Response {
	response := Response{httpVersion: request.httpVersion}

	if !request.method.isValid() {
		response.statusCode = BadRequest
		response.content = fmt.Sprintf("Invalid HTTP method %s used", request.method)
		return response
	}

	resource, err := filepath.Abs(filepath.Join(path, request.resource))
	if err != nil {
		response.statusCode = BadRequest
		response.content = fmt.Sprintf("Invalid resource. %v", err)
		return response
	}
	// TODO Explore other ways of securing the resources on the server
	if !strings.Contains(filepath.Dir(resource), path) {
		response.statusCode = BadRequest
		response.content = fmt.Sprintf("Attempting to access file outside of allowed directory.")
		return response
	}
	if _, e := os.Stat(resource); os.IsNotExist(e) {
		response.statusCode = NotFound
		response.content = fmt.Sprintf("File %s not found", resource)
		return response
	}

	data, err := ioutil.ReadFile(resource)

	if err != nil {
		response.statusCode = InternalServerError
		response.content = "Something went wrong in the server!"
		return response
	}

	// TODO Allow users to configure server name and set that as a response header
	response.statusCode = OK
	response.content = string(data)
	response.headers = make(map[string]string)
	response.headers["Content-Type"] = "text/plain"
	response.headers["Content-Length"] = strconv.Itoa(len(data))

	return response
}

func readRequest(scanner *bufio.Scanner, path string) (request Request, err error) {
	headers := make(map[string]string)
	var method Method
	var resource string
	var httpVersion string

	// First line is the start line
	if scanner.Scan() {
		startLine := scanner.Text()
		method, resource, httpVersion, err = parseStartLine(startLine)
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
