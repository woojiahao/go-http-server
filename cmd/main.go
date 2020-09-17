package main

import (
	"github.com/woojiahao/go-http-server/internal/server"
)

func main() {
	server := server.Create(8000)
	server.Start()
}
