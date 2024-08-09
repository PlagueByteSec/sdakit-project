package lib

import (
	"os"
)

type FileStreams struct {
	Ipv4AddrStream  *os.File
	Ipv6AddrStream  *os.File
	SubdomainStream *os.File
}

func OpenOutputFileStreams(params Params) (*FileStreams, error) {
	ipv4AddrStream, err := os.OpenFile(params.FilePathIPv4, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	ipv6AddrStream, err := os.OpenFile(params.FilePathIPv6, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		ipv4AddrStream.Close()
		return nil, err
	}
	subdomainStream, err := os.OpenFile(params.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		ipv4AddrStream.Close()
		ipv6AddrStream.Close()
		return nil, err
	}
	return &FileStreams{
		Ipv4AddrStream:  ipv4AddrStream,
		Ipv6AddrStream:  ipv6AddrStream,
		SubdomainStream: subdomainStream,
	}, nil
}

func WriteOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return nil
}
