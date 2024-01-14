// Package messages contains messages used in BGP protocol.
package messages

// Message is an abstraction for all BGP messages.
type Message interface {
	String() string
}
