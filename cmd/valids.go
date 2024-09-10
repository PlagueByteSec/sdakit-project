package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	utils "github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/streams"
)

func RDnsFromFile(args *shared.Args) {
	ipFileStream := streams.RoFileStreamInit(args.RDnsLookupFilePath)
	scanner := bufio.NewScanner(ipFileStream)
	fmt.Fprintln(shared.GStdout)
	shared.GStdout.Flush()
	for scanner.Scan() {
		entry := scanner.Text()
		requests.SetDnsEnumType()
		// Perform RDNS lookup against the current IP address
		requests.DnsLookups(shared.GDnsResolver, shared.DnsLookupOptions{
			IpAddress: net.ParseIP(entry),
			Subdomain: "",
		})
		fmt.Fprintf(shared.GStdout, "[+] %s\n", entry)
		for idx := 0; idx < len(shared.GDnsResults); idx++ {
			fmt.Fprintf(shared.GStdout, " | %s\n", shared.GDnsResults[idx])
		}
		shared.GStdout.Flush()
	}
	streams.ScannerCheckError(scanner)
	os.Exit(0)
}

func PingFromFile(args *shared.Args) {
	// Read subdomains from file ping each
	fileStream := streams.RoFileStreamInit(args.PingSubdomainsFile)
	scanner := bufio.NewScanner(fileStream)
	fmt.Fprintln(shared.GStdout)
	shared.GStdout.Flush()
	var (
		pingSuccess int
		pingFailed  int
	)
	for scanner.Scan() {
		entry := scanner.Text()
		if err := requests.PingSubdomain(entry, args.PingCount); err != nil {
			fmt.Printf("[-] %s: PING FAILED\n", entry)
			pingFailed++
			continue
		}
		fmt.Printf("[+] %s: PING SUCCEED\n", entry)
		pingSuccess++
	}
	streams.ScannerCheckError(scanner)
	fmt.Printf("\n[*] Summary: %d succeed, %d failed\n", pingSuccess, pingFailed)
	os.Exit(0)
}

func AnalyseHttpHeaderSingle(args *shared.Args) {
	httpClient, err := requests.HttpClientInit(args)
	if err != nil {
		utils.SentinelPanic(err)
	}
	if strings.HasPrefix(args.Subdomain, "http://") {
		args.Subdomain = strings.TrimPrefix(args.Subdomain, "http://")
	} else if strings.HasPrefix(args.Subdomain, "https://") {
		args.Subdomain = strings.TrimPrefix(args.Subdomain, "https://")
	}
	results := requests.AnalyseHttpHeader(httpClient, args.Subdomain, args.HttpRequestMethod)
	fmt.Printf("[*] Header Analysis Results For: %s\n", args.Subdomain)
	if results == "" {
		fmt.Printf("[-] Nothing to see here..\n\n")
	} else {
		fmt.Println(results)
	}
	os.Exit(0)
}
