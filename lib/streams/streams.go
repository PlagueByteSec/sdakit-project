package streams

import (
	"Sentinel/lib/shared"
	"Sentinel/lib/utils"
	"errors"
	"os"
)

func OpenOutputFileStreams(streams *shared.FileStreams, paths *shared.FilePaths) error {
	/*
		Open separate file streams for each category of output files. The categories
		are divided into IPv4 addresses, IPv6 addresses, and subdomains.
	*/
	var err error
	streams.Ipv4AddrStream, err = os.OpenFile(
		paths.FilePathIPv4,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		shared.DefaultPermission,
	)
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	streams.Ipv6AddrStream, err = os.OpenFile(
		paths.FilePathIPv6,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		shared.DefaultPermission,
	)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		shared.Glogger.Println(err)
		return err
	}
	streams.SubdomainStream, err = os.OpenFile(
		paths.FilePathSubdomain,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		shared.DefaultPermission,
	)
	if err != nil {
		streams.Ipv4AddrStream.Close()
		streams.Ipv6AddrStream.Close()
		shared.Glogger.Println(err)
		return err
	}
	return nil
}

func WriteOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	return nil
}

func CloseOutputFileStreams(streams *shared.FileStreams) {
	streams.Ipv4AddrStream.Close()
	streams.Ipv6AddrStream.Close()
	streams.SubdomainStream.Close()
}

func WordlistStreamInit(args *shared.Args) (*os.File, int) {
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find wordlist: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	lineCount, err := utils.FileCountLines(args.WordlistPath)
	if err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Failed to count lines of " + args.WordlistPath,
			ExitError:   err,
		})
	}
	wordlistStream, err := os.Open(args.WordlistPath)
	if err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	return wordlistStream, lineCount
}

func IpFileStreamInit(args *shared.Args) *os.File {
	if _, err := os.Stat(args.RDnsLookupFilePath); errors.Is(err, os.ErrNotExist) {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find IP list: " + args.RDnsLookupFilePath,
			ExitError:   err,
		})
	}
	n, err := utils.FileCountLines(args.RDnsLookupFilePath)
	if n == 0 {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Could not process an empty file: " + args.RDnsLookupFilePath,
			ExitError:   nil,
		})
	}
	ipListStream, err := os.Open(args.RDnsLookupFilePath)
	if err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.RDnsLookupFilePath,
			ExitError:   err,
		})
	}
	return ipListStream
}
