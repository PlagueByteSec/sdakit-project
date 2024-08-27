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

func fileValidate(filePath string) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find: " + filePath,
			ExitError:   err,
		})
	}
}

func openFileStreamSingle(filePath string) *os.File {
	fileStream, err := os.Open(filePath)
	if err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + filePath,
			ExitError:   err,
		})
	}
	return fileStream
}

func WordlistStreamInit(args *shared.Args) (*os.File, int) {
	fileValidate(args.WordlistPath)
	lineCount, err := utils.FileCountLines(args.WordlistPath)
	if err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Failed to count lines of " + args.WordlistPath,
			ExitError:   err,
		})
	}
	wordlistStream := openFileStreamSingle(args.WordlistPath)
	return wordlistStream, lineCount
}

func RoFileStreamInit(filePath string) *os.File {
	fileValidate(filePath)
	n, err := utils.FileCountLines(filePath)
	if n == 0 {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Could not process an empty file: " + filePath,
			ExitError:   nil,
		})
	}
	ipListStream := openFileStreamSingle(filePath)
	return ipListStream
}
