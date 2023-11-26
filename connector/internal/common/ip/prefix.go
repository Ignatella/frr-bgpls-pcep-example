package ip

import "net"

type Prefix struct {
	Prefix       net.IP
	PrefixLength uint8
}
