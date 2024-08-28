package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/internal/streams"
)

func RDnsFromFile(args *shared.Args) {
	ipFileStream := streams.RoFileStreamInit(args.RDnsLookupFilePath)
	scanner := bufio.NewScanner(ipFileStream)
	fmt.Fprintln(shared.GStdout)
	shared.GStdout.Flush()
	for scanner.Scan() {
		entry := scanner.Text()
		shared.GDnsResolver = requests.DnsResolverInit(false)
		if shared.CustomDnsServer != "" {
			// Use custom DNS server address
			shared.GDnsResolver = requests.DnsResolverInit(true)
		}
		// Perform DNS lookup against the current subdomain
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
