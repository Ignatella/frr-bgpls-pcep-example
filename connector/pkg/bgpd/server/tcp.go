package server

import (
	"connector/pkg/bgpd/events"
	"connector/pkg/bgpd/router"
	"fmt"
	"log"
	"net"
)

// Server represents a TCP Server. It listens on Host:Port
type Server struct {
	Host string
	Port string
}

// Config represents Daemon configuration. For now, it contains only Host and Port.
type Config struct {
	Host string
	Port string
}

// NewTCPServer creates and initializes a new TCP Server.
func NewTCPServer(config *Config) *Server {
	return &Server{
		Host: config.Host,
		Port: config.Port,
	}
}

// Run starts the TCP server.
// It is responsible for listening and accepting new connections.
// When a new connection is accepted, it creates a new "Router" (what represents a connection).
// From this moment, the "Router" is responsible for handling incoming messages.
func (server *Server) Run(bgpdEventCh chan events.BGPdEvent) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.Host, server.Port))
	if err != nil {
		log.Fatal(err)
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	log.Printf("Listening on %s:%s\n", server.Host, server.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		router := router.NewRouter(conn)

		go router.RunControlThread(bgpdEventCh)
		go router.HandleRequest(bgpdEventCh)
	}
}
