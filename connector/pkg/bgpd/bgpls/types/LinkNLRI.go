package types

type LinkNLRI struct {
	Endpoints []LinkEndpoint
}

type LinkEndpoint struct {
	AS             uint16
	RouterId       string
	NeighborPrefix string
}
