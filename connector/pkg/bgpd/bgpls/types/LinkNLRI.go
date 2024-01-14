package types

// LinkNLRI is a BGP-LS Link NLRI. It contains information about a link. It is defined in [RFC 7752].
//
// [RFC 7752]: https://datatracker.ietf.org/doc/html/rfc7752#section-3.3.2
type LinkNLRI struct {
	Endpoints []LinkEndpoint
}

// LinkEndpoint represents one of two endpoints of a link.
type LinkEndpoint struct {
	AS             uint16
	RouterId       string
	NeighborPrefix string
}
