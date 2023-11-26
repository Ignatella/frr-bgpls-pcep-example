package bgprouter

import (
	"connector/internal/filereader"
	"connector/pkg/bgpd/bgp/types"
	"log"
)

const (
	openMessagePath        = "bgpOpenMessage.txt"
	keepAliveMessagePath   = "bgpKeepAliveMessage.txt"
	emptyUpdateMessagePath = "bgpEmptyUpdate.txt"
)

func (router *Router) RunControlThread(bgpdEventCh chan types.BGPdEvent) {
	for {
		select {
		case <-router.running:
			return
		case event := <-router.routerEventCh:
			switch event.Type {
			case types.NewRouter:
				router := event.Data.(Router)
				log.Printf("New router: %s\n", &router)
			case types.OpenMessageReceived:
				router.writeMessageFromFile(openMessagePath)
			case types.KeepAliveMessageReceived:
				router.writeMessageFromFile(keepAliveMessagePath)
			case types.UpdateMessageReceived:
				router.writeMessageFromFile(emptyUpdateMessagePath)
			case types.Quit, types.NotificationMessageReceived:
				close(router.running)
				return
			default:
				log.Fatalf("Unknown event type: %d\n", event.Type)
			}
		}
	}
}

func (router *Router) writeMessageFromFile(path string) {
	messageBytes, err := filereader.ReadMessage(path)
	if err != nil {
		close(router.running)
		return
	}

	_, err = router.conn.Write(messageBytes)
	if err != nil {
		close(router.running)
		return
	}
}
