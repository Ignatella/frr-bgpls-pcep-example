package router

import (
	"bufio"
	"bytes"
	"connector/internal/filereader"
	"connector/internal/printer"
	"connector/pkg/bgp"
	"connector/pkg/bgp/messages"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
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

func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	log.Printf("Listening on %s:%s", server.host, server.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn: conn,
		}

		go client.handleRequest()
	}
}

const bgpMarker = "ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff"

func (client *Client) handleRequest() {
	log.Printf("New connection from %s", client.conn.RemoteAddr().String())

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(client.conn)

	reader := bufio.NewReader(client.conn)
	buffer := make([]byte, 1024)

	for {
		// read from connection
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("Error reading from %s: %s", client.conn.RemoteAddr().String(), err)
			return
		}

		// print bytes
		log.Println("Bytes incoming: ")
		printer.Hexdump(buffer[:n])

		bgpMarker, _ := hex.DecodeString(strings.Replace(bgpMarker, " ", "", -1))

		for i := 0; i < len(buffer[:n]); {
			if !bytes.Equal(buffer[i:i+16], bgpMarker) {
				log.Printf("Invalid BGP marker met\n")
				break
			}

			messageLengthBytes := buffer[i+16 : i+18]
			messageLength := int(uint16(messageLengthBytes[0])<<8 | uint16(messageLengthBytes[1]))

			printer.Hexdump(buffer[i : i+messageLength])

			// parse message
			message, err := bgp.ParseMessage(buffer[i : i+messageLength])
			if err != nil {
				log.Printf("Error parsing BGP message: %s", err)
				return
			}

			// prepare response message
			if _, ok := message.(*messages.OpenMessage); ok {
				messageBytes, err := filereader.ReadMessage("bgpOpenMessage.txt")
				if err != nil {
					log.Printf("Error reading message file: %s", err)
					return
				}

				_, err = client.conn.Write(messageBytes)
				if err != nil {
					log.Printf("Error writing to %s: %s", client.conn.RemoteAddr().String(), err)
					return
				}
			}
			if _, ok := message.(*messages.KeepAliveMessage); ok {
				messageBytes, err := filereader.ReadMessage("bgpKeepAliveMessage.txt")
				if err != nil {
					log.Printf("Error reading message file: %s", err)
					return
				}

				_, err = client.conn.Write(messageBytes)
				if err != nil {
					log.Printf("Error writing to %s: %s", client.conn.RemoteAddr().String(), err)
					return
				}
			}
			if _, ok := message.(*messages.UpdateMessage); ok {
				messageBytes, err := filereader.ReadMessage("bgpEmptyUpdate.txt")
				if err != nil {
					log.Printf("Error reading message file: %s", err)
					return
				}

				_, err = client.conn.Write(messageBytes)
				if err != nil {
					log.Printf("Error writing to %s: %s", client.conn.RemoteAddr().String(), err)
					return
				}
			}

			i += messageLength
		}
	}
}
