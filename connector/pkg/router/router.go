package router

import (
	"bufio"
	"connector/pkg/bgp"
	"connector/pkg/bgp/messages"
	"connector/pkg/printer"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
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

func (client *Client) handleRequest() {
	log.Printf("New connection from %s", client.conn.RemoteAddr().String())

	reader := bufio.NewReader(client.conn)
	// create byte buffer
	buffer := make([]byte, 1024)

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(client.conn)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("Error reading from %s: %s", client.conn.RemoteAddr().String(), err)
			return
		}

		// print bytes
		log.Println("Bytes incoming: ")
		printer.Hexdump(buffer[:n])

		// parse message
		message, err := bgp.ParseMessage(buffer[:n])
		if err != nil {
			log.Printf("Error parsing BGP message: %s", err)
			return
		}

		// prepare message
		if _, ok := message.(*messages.OpenMessage); ok {
			// read file
			assetPath := filepath.Join("assets", "bgpOpenMessage.txt")

			content, err := os.ReadFile(assetPath)
			if err != nil {
				fmt.Println("Error reading binary file:", err)
				return
			}

			re := regexp.MustCompile(`(?mU)^[0-9a-fA-F]+:(.+)(//.+)?$`)

			// Find all matches in the input
			matches := re.FindAllStringSubmatch(string(content), -1)

			var messageBytes []byte

			for _, match := range matches {
				byteString := strings.ReplaceAll(match[1], " ", "")
				bytes, err := hex.DecodeString(byteString)
				if err != nil {
					fmt.Println("Error decoding hex string:", err)
					return
				}
				messageBytes = append(messageBytes, bytes...)

			}

			_, err = client.conn.Write(messageBytes)
			if err != nil {
				log.Printf("Error writing to %s: %s", client.conn.RemoteAddr().String(), err)
				return
			}
		}

		if _, ok := message.(*messages.KeepAliveMessage); ok {
			// read file
			assetPath := filepath.Join("assets", "bgpKeepAliveMessage.txt")

			content, err := os.ReadFile(assetPath)
			if err != nil {
				fmt.Println("Error reading binary file:", err)
				return
			}

			re := regexp.MustCompile(`(?mU)^[0-9a-fA-F]+:(.+)(//.+)?$`)

			// Find all matches in the input
			matches := re.FindAllStringSubmatch(string(content), -1)

			var messageBytes []byte

			for _, match := range matches {
				byteString := strings.ReplaceAll(match[1], " ", "")
				bytes, err := hex.DecodeString(byteString)
				if err != nil {
					fmt.Println("Error decoding hex string:", err)
					return
				}
				messageBytes = append(messageBytes, bytes...)

			}

			_, err = client.conn.Write(messageBytes)
			if err != nil {
				log.Printf("Error writing to %s: %s", client.conn.RemoteAddr().String(), err)
				return
			}
		}
	}
}
