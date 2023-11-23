package filereader

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ReadMessage(filename string) ([]byte, error) {
	assetPath := filepath.Join("assets", filename)

	content, err := os.ReadFile(assetPath)

	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?mU)^[0-9a-fA-F]+:(.+)(//.+)?$`)

	matches := re.FindAllStringSubmatch(string(content), -1)

	var messageBytes []byte

	for _, match := range matches {
		byteString := strings.ReplaceAll(match[1], " ", "")
		bytes, err := hex.DecodeString(byteString)
		if err != nil {
			return nil, err
		}
		messageBytes = append(messageBytes, bytes...)
	}

	return messageBytes, nil
}
