package server

import (
	"connector/pkg/bgpd/bgpRouter"
	"connector/pkg/bgpd/events"
	"fmt"
	"log"
	"net"
)

type Server struct {
	host string
	port string
}

type Client struct {
	conn net.Conn
}

type Config struct {
	Host string
	Port string
}

func NewTCPServer(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run(bgpdEventCh chan events.BGPdEvent) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	log.Printf("Listening on %s:%s\n", server.host, server.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		router := bgprouter.NewRouter(conn)

		go router.RunControlThread(bgpdEventCh)
		go router.HandleRequest(bgpdEventCh)
	}
}
