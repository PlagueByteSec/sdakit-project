package utils

import (
	"Sentinel/lib/shared"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func DefaultOutputName(hostname string, fileExtension shared.FileExtension) string {
	var extension string
	switch fileExtension {
	case shared.TXT:
		extension = "txt"
	case shared.JSON:
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
		err := os.MkdirAll(dirname, shared.DefaultPermission)
		if err != nil {
			shared.Glogger.Println(err)
			return errors.New("unable to create output directory: " + dirname)
		}
	}
	return nil
}

func outputFileAlreadyExist(outputFilePath string) bool {
	if _, err := os.Stat(outputFilePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func clearFileContent(outputFilePath string) error {
	stream, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, shared.DefaultPermission)
	if stream != nil {
		return err
	}
	stream.Close()
	return nil
}

func cleanExistingOutputFiles(outputFiles []string) {
	// Recreate existing output files to prevent saving duplicate entries
	for idx := 0; idx < len(outputFiles); idx++ {
		file := outputFiles[idx]
		if outputFileAlreadyExist(file) {
			clearFileContent(file)
		}
	}
}

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
		VerbosePrint("[*] New output directory path set: %s\n", args.NewOutputDirPath)
	}
	if err := CreateOutputDir(args.NewOutputDirPath); err != nil {
		shared.Glogger.Println(err)
		return nil, err
	}
	var (
		filePathSubdomain string
		filePathIPv4      string
		filePathIPv6      string
		filePathJSON      string
		extension         shared.FileExtension = shared.TXT
		outputFiles       []string
	)
	if args.OutFileSubdoms == "defaultSd" {
		filePathSubdomain = filepath.Join(
			args.NewOutputDirPath,
			"Subdomains-"+DefaultOutputName(args.Domain, extension),
		)
	} else {
		filePathSubdomain = args.OutFileSubdoms
	}
	outputFiles = append(outputFiles, filePathSubdomain)
	if args.OutFileIPv4 == "defaultV4" {
		filePathIPv4 = filepath.Join(
			args.NewOutputDirPath,
			"IPv4Addresses-"+DefaultOutputName(args.Domain, extension),
		)
	} else {
		filePathIPv4 = args.OutFileIPv4
	}
	outputFiles = append(outputFiles, filePathIPv4)
	if args.OutFileIPv6 == "defaultV6" {
		filePathIPv6 = filepath.Join(
			args.NewOutputDirPath,
			"IPv6Addresses-"+DefaultOutputName(args.Domain, extension),
		)
	} else {
		filePathIPv6 = args.OutFileIPv6
	}
	outputFiles = append(outputFiles, filePathIPv6)
	if args.OutFileJSON == "defaultJSON" {
		extension = shared.JSON
		filePathJSON = filepath.Join(
			args.NewOutputDirPath,
			"Summary-"+DefaultOutputName(args.Domain, extension),
		)
	} else {
		filePathJSON = args.OutFileJSON
	}
	outputFiles = append(outputFiles, filePathJSON)
	if args.RDnsLookupFilePath == "" {
		cleanExistingOutputFiles(outputFiles)
	}
	return &shared.FilePaths{
		FilePathSubdomain: filePathSubdomain,
		FilePathIPv4:      filePathIPv4,
		FilePathIPv6:      filePathIPv6,
		FilePathJSON:      filePathJSON,
	}, nil
}

func FileCountLines(filePath string) (int, error) {
	/*
		FileCountLines counts the number of lines in a file by reading the content
		in 32 KB chunks and counting the newline characters.
	*/
	stream, err := os.Open(filePath)
	if err != nil {
		shared.Glogger.Println(err)
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
			shared.Glogger.Println(err)
			return counter, err
		}
	}
}

func ScannerCheckError(scanner *bufio.Scanner) {
	// Handle errors for wordlist scanner
	if err := scanner.Err(); err != nil {
		shared.Glogger.Println(err)
		SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Scanner failed",
			ExitError:   err,
		})
	}
}
