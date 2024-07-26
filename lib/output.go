package lib

import (
	"os"
)

func WriteOutput(filePath string, content string) {
	stream, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		GetPanic("ERROR: %s\n", err)
	}
	defer stream.Close()
	if _, err = stream.WriteString(content); err != nil {
		GetPanic("ERROR: %s\n", err)
	}
}
