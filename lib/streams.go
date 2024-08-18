package lib

import (
	"os"
)

type FileStreams struct {
	Ipv4AddrStream  *os.File
	Ipv6AddrStream  *os.File
	SubdomainStream *os.File
}

func (streams *FileStreams) OpenOutputFileStreams(paths *FilePaths) error {
	var err error
	streams.Ipv4AddrStream, err = os.OpenFile(paths.FilePathIPv4, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		Logger.Println(err)
		return err
	}
	streams.Ipv6AddrStream, err = os.OpenFile(paths.FilePathIPv6, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		Logger.Println(err)
		return err
	}
	streams.SubdomainStream, err = os.OpenFile(paths.FilePathSubdomain, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		streams.Ipv6AddrStream.Close()
		Logger.Println(err)
		return err
	}
	return nil
}

func WriteOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		Logger.Println(err)
		return err
	}
	return nil
}

func (streams *FileStreams) CloseOutputFileStreams() {
	streams.Ipv4AddrStream.Close()
	streams.Ipv6AddrStream.Close()
	streams.SubdomainStream.Close()
}
