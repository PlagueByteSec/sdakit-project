package lib

import (
	"errors"
	"os"
)

func OpenOutputFileStreamIPv4(params Params) (*os.File, error) {
	return openOutputFileStream(params.FilePathIPv4)
}

func WriteOutputFileStreamIPv4(stream *os.File, param Params) error {
	return writeOutputFileStream(stream, param.FileContentIPv4)
}

func OpenOutputFileStreamIPv6(params Params) (*os.File, error) {
	return openOutputFileStream(params.FilePathIPv6)
}

func WriteOutputFileStreamIPv6(stream *os.File, param Params) error {
	return writeOutputFileStream(stream, param.FileContentIPv6)
}

func OpenOutputFileStreamDomains(params Params) (*os.File, error) {
	return openOutputFileStream(params.FilePath)
}

func WriteOutputFileStreamDomains(stream *os.File, params Params) error {
	return writeOutputFileStream(stream, params.FileContent)
}

func openOutputFileStream(filePath string) (*os.File, error) {
	stream, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, errors.New("failed to open output file stream")
	}
	return stream, nil
}

func writeOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		return errors.New("output write operation failed")
	}
	return nil
}
