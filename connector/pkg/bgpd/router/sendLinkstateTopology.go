package router

import (
	"connector/internal/filereader"
	"connector/internal/printer"
	bgplsTypes "connector/pkg/bgpd/bgpls/types"
	"golang.org/x/exp/slices"
	"net"
)

// sendBGPLSLinkNLRI sends BGP-LS Node NLRI to the BGP-LS peer.
// It loads a template for a message from a file, fills it with the data and sends to the connection.
func (router *Router) sendBGPLSLinkNLRI(linkNlris []bgplsTypes.LinkNLRI) error {
	for _, nlri := range linkNlris {
		messageBytes, err := filereader.ReadMessage("bgplsLinkNlri.txt")
		if err != nil {
			return err
		}

		// set next hop
		messageBytes[0x1f] = 0xac
		messageBytes[0x20] = 0x14
		messageBytes[0x21] = 0x00
		messageBytes[0x22] = 0x01

		//// local node
		localRouterId := net.ParseIP(nlri.Endpoints[0].RouterId).To4()

		// as number
		messageBytes[0x3b] = 0x00
		messageBytes[0x3c] = 0x01

		// igp router id
		messageBytes[0x4e] = localRouterId[3]

		//// remote node
		remoteRouterId := net.ParseIP(nlri.Endpoints[1].RouterId).To4()

		// as number
		messageBytes[0x59] = 0x00
		messageBytes[0x5a] = 0x01

		// igp router id
		messageBytes[0x6c] = remoteRouterId[3]

		//// link descriptors

		// local interface address
		localInterfaceAddress := net.ParseIP(nlri.Endpoints[0].NeighborPrefix).To4()
		messageBytes = slices.Replace(messageBytes, 0x71, 0x75, []byte(localInterfaceAddress)...)

		// remote interface address
		remoteInterfaceAddress := net.ParseIP(nlri.Endpoints[1].NeighborPrefix).To4()
		messageBytes = slices.Replace(messageBytes, 0x79, 0x7d, []byte(remoteInterfaceAddress)...)

		//// link attributes
		// local node router id
		messageBytes = slices.Replace(messageBytes, 0x92, 0x96, []byte(localRouterId)...)

		// remote node router id
		messageBytes = slices.Replace(messageBytes, 0x9a, 0x9e, []byte(remoteRouterId)...)

		// adjacency sid tlv
		sid := 24000 + uint32(localRouterId[3])*10 + uint32(remoteRouterId[3])
		sidBytes := []byte{byte(sid >> 16), byte(sid >> 8), byte(sid)}
		messageBytes = slices.Replace(messageBytes, 0xf1, 0xf4, sidBytes...)

		printer.Hexdump(messageBytes)

		_, err = router.conn.Write(messageBytes)

		if err != nil {
			return err
		}
	}

	return nil
}

// sendBGPLSPrefixNLRI sends BGP-LS Prefix NLRI to the BGP-LS peer.
// Similar to sendBGPLSLinkNLRI.
func (router *Router) sendBGPLSPrefixNLRI(prefixNlris []bgplsTypes.PrefixNLRI) error {
	for _, nlri := range prefixNlris {

		messageBytes, err := filereader.ReadMessage("bgplsPrefixNlri.txt")
		if err != nil {
			return err
		}

		// set next hop
		messageBytes[0x1f] = 0xac
		messageBytes[0x20] = 0x14
		messageBytes[0x21] = 0x00
		messageBytes[0x22] = 0x01

		// as number
		messageBytes[0x3b] = 0x00
		messageBytes[0x3c] = 0x01

		routerId := net.ParseIP(nlri.RouterId).To4()

		// igp router id
		messageBytes[0x4e] = routerId[3]

		// prefix
		messageBytes = slices.Replace(messageBytes, 0x54, 0x58, []byte(routerId)...)

		// mpls label
		messageBytes[0x7b] = routerId[3]

		// router identifier of prefix originator
		messageBytes = slices.Replace(messageBytes, 0x85, 0x89, []byte(routerId)...)

		printer.Hexdump(messageBytes)

		_, err = router.conn.Write(messageBytes)

		if err != nil {
			return err
		}
	}

	return nil
}

// sendBGPLSNodeNLRI sends BGP-LS Node NLRI to the BGP-LS peer.
// Similar to sendBGPLSLinkNLRI.
func (router *Router) sendBGPLSNodeNLRI(nodeNlris []bgplsTypes.NodeNLRI) error {

	for _, nlri := range nodeNlris {

		messageBytes, err := filereader.ReadMessage("bgplsNodeNlri.txt")
		if err != nil {
			return err
		}

		// set next hop
		messageBytes[31] = 0xac
		messageBytes[32] = 0x14
		messageBytes[33] = 0x00
		messageBytes[34] = 0x01

		// as number
		messageBytes[59] = 0x00
		messageBytes[60] = 0x01

		routerId := net.ParseIP(nlri.RouterId).To4()

		// igp id
		messageBytes[0x4e] = routerId[3]

		// router id
		messageBytes = slices.Replace(messageBytes, 112, 116, []byte(routerId)...)

		// router name length
		routerNameLength := uint16(messageBytes[98])<<8 | uint16(messageBytes[99])
		lenDiff := len(nlri.Hostname) - int(routerNameLength)
		// update router name
		messageBytes = slices.Replace(messageBytes, 100, 100+int(routerNameLength), []byte(nlri.Hostname)...)

		// update length
		routerNameLength += uint16(lenDiff)
		messageBytes[98] = byte(routerNameLength >> 8)
		messageBytes[99] = byte(routerNameLength)

		// bgpls attributes length
		messageBytes[0x5f] += byte(lenDiff)

		// total path attribute length
		messageBytes[0x15] = byte(((uint16(messageBytes[0x15])<<8 | uint16(messageBytes[0x16])) + uint16(lenDiff)) >> 8)
		messageBytes[0x16] = byte((uint16(messageBytes[0x15])<<8 | uint16(messageBytes[0x16])) + uint16(lenDiff))

		// total message length
		messageBytes[0x10] = byte(((uint16(messageBytes[0x10])<<8 | uint16(messageBytes[0x11])) + uint16(lenDiff)) >> 8)
		messageBytes[0x11] = byte((uint16(messageBytes[0x10])<<8 | uint16(messageBytes[0x11])) + uint16(lenDiff))

		printer.Hexdump(messageBytes)

		_, err = router.conn.Write(messageBytes)

		if err != nil {
			return err
		}
	}

	return nil
}
