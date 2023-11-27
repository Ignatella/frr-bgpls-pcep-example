package bgprouter

import (
	"connector/internal/filereader"
	"connector/internal/printer"
	bgplsTypes "connector/pkg/bgpd/bgpls/types"
	"golang.org/x/exp/slices"
	"net"
)

func (router *Router) sendBGPLSTopology(nodeNlri []bgplsTypes.NodeNLRI) error {

	router.writeMessageFromFile(keepAliveMessagePath)

	for _, nlri := range nodeNlri {

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
	}

	filereader.ReadMessage(keepAliveMessagePath)

	return nil
}
