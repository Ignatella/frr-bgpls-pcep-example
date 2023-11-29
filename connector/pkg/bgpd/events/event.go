package events

//  bgpd events

const (
	NewRouter = iota
)

// router events

const (
	OpenMessageReceived = 100 + iota
	UpdateMessageReceived
	KeepAliveMessageReceived
	NotificationMessageReceived
	SendBGPLSNodeNLRI
	SendBGPLSPrefixNLRI
	SendBGPLSLinkNLRI
	Quit
)

type EventType int

type BGPdEvent struct {
	Type EventType
	Data interface{}
}

type RouterEvent struct {
	Type EventType
	Data interface{}
}
