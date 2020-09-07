package main

import (
	"fmt"
	"github.com/woojiahao/go-http-server/internal/server"
)

func main() {
	server := server.Create(8000)

	for {
		conn, err := server.Ln.Accept()
		if err != nil {
			panic(err)
		}

		client := server.AddConn(conn)
		fmt.Printf("New connection made to client %s\n", client.ID)
		go server.HandleConn(client)
	}
}
