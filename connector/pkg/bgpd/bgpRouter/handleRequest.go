package bgprouter

import (
	"bufio"
	"bytes"
	"connector/internal/printer"
	"connector/pkg/bgpd/bgp"
	"connector/pkg/bgpd/bgp/messages"
	"connector/pkg/bgpd/bgp/types"
	"log"
	"net"
)

func (router *Router) HandleRequest(eventCh chan types.BGPdEvent) {
	log.Printf("New connection from %s", router.conn.RemoteAddr().String())

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(router.conn)

	reader := bufio.NewReader(router.conn)
	buffer := make([]byte, 1024)

	acceptConnection := func() {
		// read from connection
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("Error reading from %s: %s", router.conn.RemoteAddr().String(), err)
			router.Exit()
			return
		}

		// print received bytes
		log.Printf("Received %d bytes from %s", n, router.conn.RemoteAddr().String())

		// parse received bytes
		for i := 0; i < len(buffer[:n]); {
			// check if BGP marker is valid
			bgpMarker := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
			if !bytes.Equal(buffer[i:i+16], bgpMarker) {
				log.Printf("Invalid BGP marker met\n")
				break
			}

			// get message length
			messageLengthBytes := buffer[i+16 : i+18]
			messageLength := int(uint16(messageLengthBytes[0])<<8 | uint16(messageLengthBytes[1]))

			// print message
			printer.Hexdump(buffer[i : i+messageLength])

			// parse message
			message, err := bgp.ParseMessage(buffer[i : i+messageLength])
			if err != nil {
				log.Printf("Error parsing BGP message: %s\n", err)
				break
			}

			// prepare response message
			if m, ok := message.(*messages.OpenMessage); ok {
				router.AS = m.AS
				router.BGPIdentifier = m.BGPIdentifier.String()
				router.Capabilities = m.Capabilities

				router.RouterEventCh <- types.RouterEvent{
					Type: types.OpenMessageReceived,
				}

				eventCh <- types.BGPdEvent{
					Type: types.NewRouter,
					Data: router,
				}
			}
			if _, ok := message.(*messages.KeepAliveMessage); ok {
				router.RouterEventCh <- types.RouterEvent{
					Type: types.KeepAliveMessageReceived,
				}
			}
			if _, ok := message.(*messages.NotificationMessage); ok {
				router.RouterEventCh <- types.RouterEvent{
					Type: types.NotificationMessageReceived,
				}
			}
			if m, ok := message.(*messages.UpdateMessage); ok {
				router.AddPrefix(m.Prefix.String())

				router.RouterEventCh <- types.RouterEvent{
					Type: types.UpdateMessageReceived,
				}

				eventCh <- types.BGPdEvent{
					Type: types.NewRouter,
					Data: router,
				}
			}

			i += messageLength
		}
	}

	for {
		select {
		case <-router.running:
			return
		default:
			acceptConnection()
		}
	}
}
