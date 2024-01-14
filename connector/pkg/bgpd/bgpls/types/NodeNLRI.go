// Package types contains types used in BGP-LS protocol.
package types

// NodeNLRI is a BGP-LS Node NLRI. It contains information about a node. It is defined in [RFC 7752].
//
// [RFC 7752]: https://datatracker.ietf.org/doc/html/rfc7752#section-3.3.1
type NodeNLRI struct {
	AS       uint16
	RouterId string
	Hostname string
}
