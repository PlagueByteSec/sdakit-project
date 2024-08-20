package utils

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
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
		err := os.MkdirAll(dirname, DefaultPermission)
		if err != nil {
			Glogger.Println(err)
			return errors.New("unable to create output directory: " + dirname)
		}
	}
	return nil
}

func GetCurrentLocalVersion() string {
	/*
		Read the version of the current local project instance. If an error
		occurs while trying to read version.txt, return n/a.
	*/
	cwd, err := os.Getwd()
	if err != nil {
		Glogger.Println(err)
		return NotAvailable
	}
	versionFilePath := filepath.Join(cwd, VersionFile)
	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		Glogger.Println(err)
		return NotAvailable
	}
	return string(content)
}

func VersionCompare(versionRepo string, versionLocal string) {
	/*
		Compare the version of the local project instance with the version
		from the GitHub repository. If the local version is lower than the repository
		version, the user is notified that updates are available.
	*/
	if versionRepo == NotAvailable ||
		versionLocal == NotAvailable || versionLocal == "" {
		return
	}
	parseRepoVersion, _ := version.NewVersion(versionRepo)
	parseLocalVersion, _ := version.NewVersion(versionLocal)
	if versionRepo != versionLocal && parseLocalVersion.LessThan(parseRepoVersion) {
		fmt.Fprintf(GStdout, "[*] An update is available! %s->%s\n", versionLocal, versionRepo)
	}
}

func InArgList(httpCode string, list []string) bool {
	/*
		An HTTP status code will be compared against those
		specified by the -e or -f flag. This function is used
		to filter the output for customization.
	*/
	for _, code := range list {
		if httpCode == code {
			return true
		}
	}
	return false
}

func EditDbEntries(args *Args) ([]string, error) {
	/*
		All endpoints will be read from db.go and formatted by replacing
		the placeholder (HOST) with the target domain. If a text file containing
		custom endpoints is specified by the -x flag, those will be
		added to the existing entries.
	*/
	entries := make([]string, 0, len(Db))
	for idx, entry := range Db {
		endpoint := strings.Replace(entry, Placeholder, args.Domain, 1)
		if args.Verbose {
			fmt.Fprintf(GStdout, "\n%d. Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
		}
		entries = append(entries, endpoint)
	}
	if args.DbExtendPath != "" {
		VerbosePrint("\n[*] Extending endpoints..")
		stream, err := os.Open(args.DbExtendPath)
		if err != nil {
			Glogger.Println(err)
			return nil, errors.New("failed to open file stream for: " + args.DbExtendPath)
		}
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		idx := 0
		for scanner.Scan() {
			entry := scanner.Text()
			if !strings.Contains(entry, Placeholder) {
				fmt.Fprintln(GStdout, "[-] Invalid pattern (HOST missing): "+entry)
				continue
			}
			endpoint := strings.Replace(entry, Placeholder, args.Domain, 1)
			VerbosePrint("\n%d. X Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
			entries = append(entries, endpoint)
			idx++
		}
	}
	VerbosePrint("\n[*] Using %d endpoints\n", len(entries))
	return entries, nil
}

func RequestIpAddresses(useCustomDnsServer bool, subdomain string) []string {
	/*
		Perform a DNS lookup for the current subdomain to get the corresponding
		IP addresses and filter out old and inactive subdomains.
	*/
	var resolver *net.Resolver
	switch useCustomDnsServer {
	case true:
		resolver = &net.Resolver{
			PreferGo: false,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", CustomDnsServer)
			},
		}
	case false:
		resolver = &net.Resolver{}
	}
	var results []string
	retryLookup, err := resolver.LookupIPAddr(context.Background(), subdomain)
	if err != nil {
		return nil
	}
	for _, ip := range retryLookup {
		results = append(results, ip.String())
	}
	return results
}

func GetIpVersion(ipAddress string) int {
	var ipVersion int
	if ip := net.ParseIP(ipAddress); ip != nil {
		if ip.To4() != nil {
			ipVersion = 4
		} else {
			ipVersion = 6
		}
	}
	return ipVersion
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
	stream, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DefaultPermission)
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

func FilePathInit(args *Args) (*FilePaths, error) {
	/*
		Build output file names for each category using default
		constructs or custom names specified by the -oS, -o4, and -o6 parameters.
	*/
	if args.NewOutputDirPath == "defaultPath" {
		args.NewOutputDirPath = OutputDir
	} else {
		VerbosePrint("[*] New output directory path set: %s\n", args.NewOutputDirPath)
	}
	if err := CreateOutputDir(args.NewOutputDirPath); err != nil {
		Glogger.Println(err)
		return nil, err
	}
	var (
		filePathSubdomain string
		filePathIPv4      string
		filePathIPv6      string
		filePathJSON      string
		extension         FileExtension = TXT
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
		extension = JSON
		filePathJSON = filepath.Join(
			args.NewOutputDirPath,
			"Summary-"+DefaultOutputName(args.Domain, extension),
		)
	} else {
		filePathJSON = args.OutFileJSON
	}
	outputFiles = append(outputFiles, filePathJSON)
	cleanExistingOutputFiles(outputFiles)
	return &FilePaths{
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
		Glogger.Println(err)
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
			Glogger.Println(err)
			return counter, err
		}
	}
}

func Evaluation(startTime time.Time, count int) {
	// Calculate the time duration and format the summary
	defer GStdout.Flush()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	var temp strings.Builder
	temp.WriteString("subdomain")
	if count != 1 {
		temp.WriteString("s")
	}
	fmt.Fprintf(GStdout, "\n[*] %d %s obtained, %d displayed\n", count, temp.String(), GDisplayCount)
	fmt.Fprintf(GStdout, "[*] Finished in %s\n", duration)
}

func SentinelExit(exitParams SentinelExitParams) {
	/*
		Read the exit settings specified in SentinelExitParams and
		adjust the behavior based on those settings.
	*/
	fmt.Fprintln(GStdout, exitParams.ExitMessage)
	if exitParams.ExitError != nil {
		errorMessage := fmt.Sprintf("Sentinel exit with an error: %s", exitParams.ExitError.Error())
		Glogger.Println(errorMessage)
		fmt.Fprintln(GStdout, errorMessage)
	}
	GStdout.Flush()
	os.Exit(exitParams.ExitCode)
}

func VerbosePrint(format string, args ...interface{}) {
	// Only print content if the -v flag is specified
	if GVerbose {
		fmt.Fprintf(GStdout, format, args...)
	}
}

func IsValidDomain(domain string) bool {
	/*
		Verify the target domain by checking it against a regex and
		performing a DNS lookup. The domain will be considered invalid
		if no IP addresses can be enumerated.
	*/
	regex := `^(?i)([a-z0-9](-?[a-z0-9])*)+(\.[a-z]{2,})+$`
	regexCompile := regexp.MustCompile(regex)
	if regexCompile.MatchString(domain) {
		ipAddrs, err := net.LookupIP(domain)
		if err != nil {
			return false
		}
		if len(ipAddrs) != 0 {
			return true
		}
	}
	return false
}

func IpManage(params Params, ip string, fileStream *FileStreams) {
	/*
		Request the IP version based on the given IP address string. A check
		is performed to verify that the address written to the output file is not
		duplicated. If successful, the address will be written to the appropriate output file.
	*/
	ipAddrVersion := GetIpVersion(ip)
	switch ipAddrVersion {
	case 4:
		params.FileContentIPv4Addrs = ip
		if !PoolContainsEntry(GPool.PoolIPv4Addresses, params.FileContentIPv4Addrs) {
			GPool.PoolIPv4Addresses = append(GPool.PoolIPv4Addresses, params.FileContentIPv4Addrs)
			err := WriteOutputFileStream(fileStream.Ipv4AddrStream, params.FileContentIPv4Addrs)
			if err != nil {
				fileStream.Ipv4AddrStream.Close()
				Glogger.Println(err)
			}
		}
		GSubdomBase.IpAddresses.IPv4 = append(
			GSubdomBase.IpAddresses.IPv4,
			net.ParseIP(ip),
		)
	case 6:
		params.FileContentIPv6Addrs = ip
		if !PoolContainsEntry(GPool.PoolIPv6Addresses, params.FileContentIPv6Addrs) {
			GPool.PoolIPv6Addresses = append(GPool.PoolIPv6Addresses, params.FileContentIPv6Addrs)
			err := WriteOutputFileStream(fileStream.Ipv6AddrStream, params.FileContentIPv6Addrs)
			if err != nil {
				fileStream.Ipv6AddrStream.Close()
				Glogger.Println(err)
			}
		}
		GSubdomBase.IpAddresses.IPv6 = append(
			GSubdomBase.IpAddresses.IPv6,
			net.ParseIP(ip),
		)
	}
}

func PrintProgress(entryCount int) {
	fmt.Fprintf(GStdout, "\rProgress::[%d/%d]", GAllCounter, entryCount)
	GStdout.Flush()
	GAllCounter++
}

func ScannerCheckError(scanner *bufio.Scanner) {
	// Handle errors for wordlist scanner
	if err := scanner.Err(); err != nil {
		Glogger.Println(err)
		SentinelExit(SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Scanner failed",
			ExitError:   err,
		})
	}
}
