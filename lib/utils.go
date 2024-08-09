package lib

import (
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
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return errors.New("unable to create output directory: " + outputDir)
		}
	}
	return nil
}

func GetCurrentLocalVersion(failHandler *VersionHandler) string {
	var (
		versionPath string
		content     []byte
		err         error
	)
	if runtime.GOOS == "windows" {
		versionPath = "..\\version.txt"
	} else if runtime.GOOS == "linux" {
		versionPath = "../version.txt"
	}
	if _, err := os.Stat(versionPath); errors.Is(err, os.ErrNotExist) {
		versionPath = "version.txt"
	}
	content, err = os.ReadFile(versionPath)
	version := string(content)
	// Mark version with n/a if reader failed
	TestVersionFail(*failHandler, &version, err)
	return version
}

func VersionCompare(versionRepo string, versionLocal string) {
	if versionRepo == Na || versionLocal == Na || versionLocal == "" {
		return
	}
	parseRepoVersion, _ := version.NewVersion(versionRepo)
	parseLocalVersion, _ := version.NewVersion(versionLocal)
	if versionRepo != versionLocal && parseLocalVersion.LessThan(parseRepoVersion) {
		fmt.Printf("[*] An update is available! %s->%s\n", versionLocal, versionRepo)
	}
}

func IsInExclude(httpCode string, list []string) bool {
	for _, code := range list {
		if httpCode == code {
			return true
		}
	}
	return false
}

func EditDbEntries(args *Args) []string {
	entries := make([]string, 0, len(Db))
	for idx, entry := range Db {
		endpoint := strings.Replace(entry, "HOST", args.Host, 1)
		if args.Verbose {
			fmt.Printf("\n%d. Entry: %s\n ===[ %s\n", idx+1, entry, endpoint)
		}
		entries = append(entries, endpoint)
	}
	if args.Verbose {
		fmt.Printf("\n[*] Using %d endpoints\n", len(entries))
	}
	return entries
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
