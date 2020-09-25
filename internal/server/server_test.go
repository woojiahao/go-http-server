package server

import (
	"bufio"
	"fmt"
	. "github.com/stretchr/testify/assert"
	"net"
	"sync"
	. "testing"
)

func connectUser(s *Server, wg *sync.WaitGroup, t *T) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(fmt.Sprintf("Error connecting to server: %v", err))
	}

	_, _ = conn.Write([]byte("GET hello\r\n"))

	status, _ := bufio.NewReader(conn).ReadString('\n')
	Equal(t, status, "-- You sent GET hello\n")
	status, _ = bufio.NewReader(conn).ReadString('\n')
	Equal(t, status, ">- ANSWER world\n")
}

func TestSingleUserMessage(t *T) {
	s := Create(8080)
	go s.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	connectUser(s, &wg, t)
	wg.Wait()
	s.Stop()
}
