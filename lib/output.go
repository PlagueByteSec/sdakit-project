package lib

import (
	"os"
	"strings"
)

func WriteOutput(filePath string, content string) {
	hasExt := strings.HasSuffix(filePath, ".txt")
	if !hasExt {
		filePath = filePath + ".txt"
	}
	stream, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		GetPanic("ERROR: %s\n", err)
	}
	defer stream.Close()
	if _, err = stream.WriteString(content + "\n"); err != nil {
		GetPanic("ERROR: %s\n", err)
	}
}
