package lib

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

func DefaultOutputName(hostname string) string {
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	outputFile := fmt.Sprintf("%s-%s.txt", formatTime, hostname)
	return outputFile
}

func CreateOutputDir() error {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(OutputDir, 0755)
		if err != nil {
			return errors.New("unable to create output directory: " + OutputDir)
		}
	}
	return nil
}

func GetCurrentLocalVersion() string {
	var (
		versionPath string
		content     []byte
		err         error
	)
	switch runtime.GOOS {
	case "windows":
		versionPath = "..\\version.txt"
	case "linux":
		versionPath = "../version.txt"
	}
	if _, err := os.Stat(versionPath); errors.Is(err, os.ErrNotExist) {
		versionPath = "version.txt"
	}
	content, err = os.ReadFile(versionPath)
	if err != nil {
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
