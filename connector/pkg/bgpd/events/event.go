// Package events contains events with help of which communication inside the daemon is happening.
package events

// BGP daemon events
const (
	NewRouter = iota // NewRouter is issued when a new connection is established.
)

// Router events
const (
	OpenMessageReceived         = 100 + iota // OpenMessageReceived is issued when an Open message is received.
	UpdateMessageReceived                    // UpdateMessageReceived is issued when an Update message is received.
	KeepAliveMessageReceived                 // KeepAliveMessageReceived is issued when a KeepAlive message is received.
	NotificationMessageReceived              // NotificationMessageReceived is issued when a Notification message is received.
	SendBGPLSNodeNLRI                        // SendBGPLSNodeNLRI is issued when a BGP-LS Node NLRI should be sent.
	SendBGPLSPrefixNLRI                      // SendBGPLSPrefixNLRI is issued when a BGP-LS Prefix NLRI should be sent.
	SendBGPLSLinkNLRI                        // SendBGPLSLinkNLRI is issued when a BGP-LS Link NLRI should be sent.
	Quit                                     // Quit is issued when a router should quit breaking the connection.
)

// EventType is a wrapper on known event types.
type EventType int

// BGPdEvent is and event with data issued by BGP daemon.
type BGPdEvent struct {
	Type EventType
	Data interface{}
}

// RouterEvent is and event with data issued by a router.
type RouterEvent struct {
	Type EventType
	Data interface{}
}
