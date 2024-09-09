package streams

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"

	utils "github.com/fhAnso/Sentinel/v1/internal/coreutils"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/pkg"
)

func ParamsSetupFiles(paramsFileSetup shared.ParamsSetupFilesBase) {
	*paramsFileSetup.FileParams = shared.Params{
		Domain:             paramsFileSetup.CliArgs.Domain,
		Subdomain:          paramsFileSetup.Subdomain,
		FileContentSubdoms: paramsFileSetup.Subdomain,
	}
	if paramsFileSetup.FilePaths == nil {
		// -dO specified, clear output file paths
		paramsFileSetup.FileParams.FilePathSubdomains = ""
		paramsFileSetup.FileParams.FilePathIPv4Addrs = ""
		paramsFileSetup.FileParams.FilePathIPv6Addrs = ""
	} else {
		// Setup output file paths
		paramsFileSetup.FileParams.FilePathSubdomains = paramsFileSetup.FilePaths.FilePathSubdomain
		paramsFileSetup.FileParams.FilePathIPv4Addrs = paramsFileSetup.FilePaths.FilePathIPv4
		paramsFileSetup.FileParams.FilePathIPv6Addrs = paramsFileSetup.FilePaths.FilePathIPv6
	}
}

func FilePathInit(args *shared.Args) (*shared.FilePaths, error) {
	/*
		Build output file names for each category using default
		constructs or custom names specified by the -oS, -o4, and -o6 parameters.
	*/
	if args.NewOutputDirPath == "defaultPath" {
		args.NewOutputDirPath = shared.OutputDir
	} else {
		utils.PrintVerbose("[*] New output directory path set: %s\n", args.NewOutputDirPath)
	}
	if err := pkg.CreateOutputDir(args.NewOutputDirPath); err != nil {
		shared.Glogger.Println(err)
		return nil, err
	}
	var (
		filePathSubdomain string
		filePathIPv4      string
		filePathIPv6      string
		filePathJSON      string
		extension         pkg.FileExtension = pkg.TXT
		outputFiles       []string
	)
	filePathSubdomain = filepath.Join(
		args.NewOutputDirPath,
		"Subdomains-"+pkg.DefaultOutputName(args.Domain, extension),
	)
	outputFiles = append(outputFiles, filePathSubdomain)
	filePathIPv4 = filepath.Join(
		args.NewOutputDirPath,
		"IPv4Addresses-"+pkg.DefaultOutputName(args.Domain, extension),
	)
	outputFiles = append(outputFiles, filePathIPv4)
	filePathIPv6 = filepath.Join(
		args.NewOutputDirPath,
		"IPv6Addresses-"+pkg.DefaultOutputName(args.Domain, extension),
	)
	outputFiles = append(outputFiles, filePathIPv6)
	extension = pkg.JSON
	filePathJSON = filepath.Join(
		args.NewOutputDirPath,
		"Summary-"+pkg.DefaultOutputName(args.Domain, extension),
	)
	outputFiles = append(outputFiles, filePathJSON)
	if args.RDnsLookupFilePath == "" {
		pkg.CleanExistingOutputFiles(outputFiles)
	}
	return &shared.FilePaths{
		FilePathSubdomain: filePathSubdomain,
		FilePathIPv4:      filePathIPv4,
		FilePathIPv6:      filePathIPv6,
		FilePathJSON:      filePathJSON,
	}, nil
}

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
	lineCount, err := pkg.FileCountLines(args.WordlistPath)
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
	n, err := pkg.FileCountLines(filePath)
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

func WriteOutputFileStream(stream *os.File, content string) error {
	_, err := stream.WriteString(content + "\n")
	if err != nil {
		return err
	}
	return nil
}

func ScannerCheckError(scanner *bufio.Scanner) {
	// Handle errors for wordlist scanner
	if err := scanner.Err(); err != nil {
		shared.Glogger.Println(err)
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Scanner failed",
			ExitError:   err,
		})
	}
}
