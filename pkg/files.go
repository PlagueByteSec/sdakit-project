package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type FileExtension int

const defaultPermission = 0755
const (
	TXT FileExtension = iota
	JSON
)

func DefaultOutputName(hostname string, fileExtension FileExtension) string {
	var extension string
	switch fileExtension {
	case TXT:
		extension = "txt"
	case JSON:
		extension = "json"
	}
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	outputFile := fmt.Sprintf("%s-%s.%s", formatTime, hostname, extension)
	return outputFile
}

func CreateOutputDir(dirname string) error {
	/*
		By default, create an output directory called output. An
		alternative name can be defined using the -nP flag.
	*/
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.MkdirAll(dirname, defaultPermission)
		if err != nil {
			return errors.New("unable to create output directory: " + dirname)
		}
	}
	return nil
}

func OutputFileAlreadyExist(outputFilePath string) bool {
	if _, err := os.Stat(outputFilePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func ClearFileContent(outputFilePath string) error {
	stream, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, defaultPermission)
	if err != nil {
		return err
	}
	stream.Close()
	return nil
}

func CleanExistingOutputFiles(outputFiles []string) {
	// Recreate existing output files to prevent saving duplicate entries
	for idx := 0; idx < len(outputFiles); idx++ {
		file := outputFiles[idx]
		if OutputFileAlreadyExist(file) {
			ClearFileContent(file)
		}
	}
}

func FileCountLines(filePath string) (int, error) {
	/*
		FileCountLines counts the number of lines in a file by reading the content
		in 32 KB chunks and counting the newline characters.
	*/
	stream, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer stream.Close()
	counter := 0
	newLine := []byte{'\n'}
	buffer := make([]byte, 32*1024)
	for {
		reader, err := stream.Read(buffer)
		if reader > 0 {
			counter += bytes.Count(buffer[:reader], newLine)
		}
		switch {
		case err == io.EOF:
			return counter, nil
		case err != nil:
			return counter, err
		}
	}
}

func LineIgnore(line string) bool {
	// Wordlists: Ensure empty lines an comments will be ignored
	trimLine := strings.TrimSpace(line)
	switch {
	case len(trimLine) == 0:
		return true
	case strings.HasPrefix(trimLine, "#") || strings.HasPrefix(trimLine, "//"):
		return true
	}
	return false
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
