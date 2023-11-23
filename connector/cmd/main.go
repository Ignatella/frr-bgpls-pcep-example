package main

import (
	"connector/pkg/router"
)

func main() {
	server := router.New(&router.Config{
		Host: "172.19.0.1",
		Port: "1790",
	})
	server.Run()
}
