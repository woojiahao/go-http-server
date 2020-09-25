package main

import (
	"github.com/woojiahao/go-http-server/internal/server"
)

func main() {
	s := server.Create(8000)
	s.Start()
}
