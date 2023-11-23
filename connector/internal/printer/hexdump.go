package printer

import "fmt"

func Hexdump(bytes []byte) {
	// 2 columns of 16 bytes
	const bytesPerLine = 16

	for i, b := range bytes {
		// address
		if i%bytesPerLine == 0 {
			fmt.Printf("%08x: ", i)
		}

		// space between 2 columns
		if i%8 == 0 {
			fmt.Printf(" ")
		}

		// print byte
		fmt.Printf("%02x ", b)

		// new line
		if i%bytesPerLine == bytesPerLine-1 {
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\n")
}
