package main

import (
	"connector/pkg/bgpd"
	"connector/pkg/bgpd/server"
)

func main() {
	bgpd.StartDaemon(&server.Config{
		Host: "0.0.0.0",
		Port: "1790",
	})
}
