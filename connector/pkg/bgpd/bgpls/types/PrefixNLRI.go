package types

// PrefixNLRI is a BGP-LS Prefix NLRI. It contains information about a prefix. It is defined in [RFC 7752].
//
// [RFC 7752]: https://datatracker.ietf.org/doc/html/rfc7752#section-3.3.3
type PrefixNLRI struct {
	AS       uint16
	RouterId string
}
