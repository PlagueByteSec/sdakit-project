package utils

import (
	"errors"
	"os"
)

type FileStreams struct {
	Ipv4AddrStream  *os.File
	Ipv6AddrStream  *os.File
	SubdomainStream *os.File
}

func (streams *FileStreams) OpenOutputFileStreams(paths *FilePaths) error {
	/*
		Open separate file streams for each category of output files. The categories
		are divided into IPv4 addresses, IPv6 addresses, and subdomains.
	*/
	var err error
	streams.Ipv4AddrStream, err = os.OpenFile(
		paths.FilePathIPv4,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		DefaultPermission,
	)
	if err != nil {
		Glogger.Println(err)
		return err
	}
	streams.Ipv6AddrStream, err = os.OpenFile(
		paths.FilePathIPv6,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		DefaultPermission,
	)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		Glogger.Println(err)
		return err
	}
	streams.SubdomainStream, err = os.OpenFile(
		paths.FilePathSubdomain,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		DefaultPermission,
	)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		streams.Ipv6AddrStream.Close()
		Glogger.Println(err)
		return err
	}
	return nil
}

func WriteOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		Glogger.Println(err)
		return err
	}
	return nil
}

func (streams *FileStreams) CloseOutputFileStreams() {
	streams.Ipv4AddrStream.Close()
	streams.Ipv6AddrStream.Close()
	streams.SubdomainStream.Close()
}

func WordlistStreamInit(args *Args) (*os.File, int) {
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find wordlist: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	lineCount, err := FileCountLines(args.WordlistPath)
	if err != nil {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Failed to count lines of " + args.WordlistPath,
			ExitError:   err,
		})
	}
	wordlistStream, err := os.Open(args.WordlistPath)
	if err != nil {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	return wordlistStream, lineCount
}

func IpFileStreamInit(args *Args) *os.File {
	if _, err := os.Stat(args.RDnsLookupFilePath); errors.Is(err, os.ErrNotExist) {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find IP list: " + args.RDnsLookupFilePath,
			ExitError:   err,
		})
	}
	ipListStream, err := os.Open(args.RDnsLookupFilePath)
	if err != nil {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.RDnsLookupFilePath,
			ExitError:   err,
		})
	}
	return ipListStream
}
