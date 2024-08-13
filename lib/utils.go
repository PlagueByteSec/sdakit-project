package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

type FilePaths struct {
	FilePathSubdomain string
	FilePathIPv4      string
	FilePathIPv6      string
}

func DefaultOutputName(hostname string) string {
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	outputFile := fmt.Sprintf("%s-%s.txt", formatTime, hostname)
	return outputFile
}

func CreateOutputDir(dirname string) error {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			Logger.Println(err)
			return errors.New("unable to create output directory: " + dirname)
		}
	}
	return nil
}

func GetCurrentLocalVersion() string {
	cwd, err := os.Getwd()
	if err != nil {
		Logger.Println(err)
		return NotAvailable
	}
	versionFilePath := filepath.Join(cwd, VersionFile)
	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		Logger.Println(err)
		return NotAvailable
	}
	return string(content)
}

func VersionCompare(versionRepo string, versionLocal string) {
	if versionRepo == NotAvailable || versionLocal == NotAvailable || versionLocal == "" {
		return
	}
	parseRepoVersion, _ := version.NewVersion(versionRepo)
	parseLocalVersion, _ := version.NewVersion(versionLocal)
	if versionRepo != versionLocal && parseLocalVersion.LessThan(parseRepoVersion) {
		fmt.Printf("[*] An update is available! %s->%s\n", versionLocal, versionRepo)
	}
}

func InArgList(httpCode string, list []string) bool {
	for _, code := range list {
		if httpCode == code {
			return true
		}
	}
	return false
}

func EditDbEntries(args *Args) ([]string, error) {
	entries := make([]string, 0, len(Db))
	for idx, entry := range Db {
		endpoint := strings.Replace(entry, Placeholder, args.Host, 1)
		if args.Verbose {
			fmt.Printf("\n%d. Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
		}
		entries = append(entries, endpoint)
	}
	if args.DbExtendPath != "" {
		fmt.Println("\n[*] Extending endpoints..")
		stream, err := os.Open(args.DbExtendPath)
		if err != nil {
			Logger.Println(err)
			return nil, errors.New("failed to open file stream for: " + args.DbExtendPath)
		}
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		idx := 1
		for scanner.Scan() {
			entry := scanner.Text()
			if !strings.Contains(entry, Placeholder) {
				fmt.Println("[-] Invalid pattern (HOST missing): " + entry)
				continue
			}
			endpoint := strings.Replace(entry, Placeholder, args.Host, 1)
			if args.Verbose {
				fmt.Printf("\n%d. X Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
			}
			entries = append(entries, endpoint)
			idx++
		}
	}
	if args.Verbose {
		fmt.Printf("\n[*] Using %d endpoints\n", len(entries))
	}
	return entries, nil
}

func RequestIpAddresses(subdomain string) []string {
	ips, err := net.LookupIP(subdomain)
	if err != nil {
		Logger.Println(err)
		return nil
	}
	var results []string
	for _, ip := range ips {
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

func FileCountLines(filePath string) (int, error) {
	stream, err := os.Open(filePath)
	if err != nil {
		Logger.Println(err)
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
			Logger.Println(err)
			return counter, err
		}
	}
}

func Evaluation(startTime time.Time, count int) {
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	var temp strings.Builder
	temp.WriteString("subdomain")
	if count != 1 {
		temp.WriteString("s")
	}
	fmt.Printf("\n\n[*] %d %s obtained, %d displayed\n", count, temp.String(), DisplayCount)
	fmt.Printf("[*] Finished in %s\n", duration)
}

func FilePathInit(args *Args) *FilePaths {
	var (
		filePathSubdomain string
		filePathIPv4      string
		filePathIPv6      string
	)
	if args.OutFile == "defaultSd" {
		filePathSubdomain = filepath.Join("output", DefaultOutputName(args.Host))
	} else {
		filePathSubdomain = args.OutFile
	}
	if args.OutFileIPv4 == "defaultV4" {
		filePathIPv4 = filepath.Join("output", "IPv4-"+DefaultOutputName(args.Host))
	} else {
		filePathIPv4 = args.OutFileIPv4
	}
	if args.OutFileIPv6 == "defaultV6" {
		filePathIPv6 = filepath.Join("output", "IPv6-"+DefaultOutputName(args.Host))
	} else {
		filePathIPv6 = args.OutFileIPv6
	}
	return &FilePaths{
		FilePathSubdomain: filePathSubdomain,
		FilePathIPv4:      filePathIPv4,
		FilePathIPv6:      filePathIPv6,
	}
}
