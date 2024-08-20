package lib

import (
	"Sentinel/lib/utils"
	"errors"
	"flag"
	"fmt"
	"os"
)

func CliParser() (utils.Args, error) {
	verbose := flag.Bool("v", false, "Verbose output")
	domain := flag.String("d", "", "Set the target domain name")
	outFile := flag.String("oS", "defaultSd", "Output file path for subdomains")
	outFileIPv4 := flag.String("o4", "defaultV4", "Output file path for IPv4 addresses")
	outFileIPv6 := flag.String("o6", "defaultV6", "Output file path for IPv6 addresses")
	outFileJSON := flag.String("oJ", "defaultJSON", "Output file path for JSON summary")
	newOutputPath := flag.String("nP", "defaultPath", "Output directory path for all results")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each subdomain")
	wordlistPath := flag.String("w", "", "Specify wordlist and direct bruteforce subdomains")
	excHttpCodes := flag.String("e", "", "Exclude HTTP codes (comma seperated)")
	filtHttpCodes := flag.String("f", "", "Filter for specific HTTP response codes (comma seperated)")
	analyzeHeader := flag.Bool("a", false, "Analyze HTTP header of each subdomain")
	portScan := flag.String("p", "", "Define port range an run scan")
	dbExtendPath := flag.String("x", "", "Extend endpoint DB with custom list")
	timeout := flag.Int("t", 5, "Specify the request timeout")
	torRoute := flag.Bool("r", false, "Enable TOR routing")
	dnsLookup := flag.Bool("dns", false, "Use wordlist (-w) and resolve subdomains by querying a DNS")
	dnsLookupCustom := flag.String("dnsC", "", "Specify custom DNS server")
	dnsLookupTimeout := flag.Int("dnsT", 500, "Specify timeout for DNS queries in ms")
	httpRequestDelay := flag.Int("rD", 500, "Specify HTTP request delay")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(Help + "\nPlease specify a domain!")
		os.Exit(0)
	}
	if *excHttpCodes != "" && !*httpCode || *filtHttpCodes != "" && !*httpCode {
		return utils.Args{}, errors.New("HTTP code filter enabled, but status codes not requested")
	}
	if !utils.IsValidDomain(*domain) {
		return utils.Args{}, errors.New("domain verification failed: " + *domain)
	}
	if !*dnsLookup && *dnsLookupCustom != "" {
		return utils.Args{}, errors.New("custom DNS address can only be set when -dns is specified")
	}
	if *dnsLookup && *wordlistPath == "" {
		return utils.Args{}, errors.New("no wordlist specified, dns method cannot be used")
	}
	args := utils.Args{
		Verbose:          *verbose,
		Domain:           *domain,
		OutFileSubdoms:   *outFile,
		OutFileIPv4:      *outFileIPv4,
		OutFileIPv6:      *outFileIPv6,
		OutFileJSON:      *outFileJSON,
		NewOutputDirPath: *newOutputPath,
		HttpCode:         *httpCode,
		WordlistPath:     *wordlistPath,
		ExcHttpCodes:     *excHttpCodes,
		FilHttpCodes:     *filtHttpCodes,
		AnalyzeHeader:    *analyzeHeader,
		PortScan:         *portScan,
		DbExtendPath:     *dbExtendPath,
		Timeout:          *timeout,
		TorRoute:         *torRoute,
		DnsLookup:        *dnsLookup,
		DnsLookupCustom:  *dnsLookupCustom,
		DnsLookupTimeout: *dnsLookupTimeout,
		HttpRequestDelay: *httpRequestDelay,
	}
	return args, nil
}
