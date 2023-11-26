package printer

import (
	"fmt"
	"log"
)

func Hexdump(bytes []byte) {

	dump := "\n"

	// 2 columns of 16 bytes
	const bytesPerLine = 16

	for i, b := range bytes {
		// address
		if i%bytesPerLine == 0 {
			dump += fmt.Sprintf("%08x: ", i)
		}

		// space between 2 columns
		if i%8 == 0 {
			dump += fmt.Sprintf(" ")
		}

		// print byte
		dump += fmt.Sprintf("%02x ", b)

		// new line
		if i%bytesPerLine == bytesPerLine-1 {
			dump += fmt.Sprintf("\n")
		}
	}

	dump += fmt.Sprintf("\n")

	log.Printf(dump)
}
