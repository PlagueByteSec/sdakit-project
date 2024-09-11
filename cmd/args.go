package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/cli"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

func CliParser() (shared.Args, error) {
	help := flag.Bool("h", false, "Display this help banner")
	verbose := flag.Bool("v", false, "Verbose output")
	domain := flag.String("d", "", "Set the target domain name")
	subdomain := flag.String("s", "", "Set the target subdomain")
	ipAddresses := flag.String("i", "", "Specify target IP address")
	newOutputPath := flag.String("nP", "defaultPath", "Output directory path for all results")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each subdomain")
	wordlistPath := flag.String("w", "", "Specify wordlist and direct bruteforce subdomains")
	excHttpCodes := flag.String("e", "", "Exclude HTTP codes (comma seperated)")
	filtHttpCodes := flag.String("f", "", "Filter for specific HTTP response codes (comma seperated)")
	analyzeHeader := flag.Bool("a", false, "Analyze HTTP header of each subdomain")
	portScan := flag.String("p", "", "Define port range an run scan")
	dbExtendPath := flag.String("x", "", "Extend endpoint DB with custom list")
	timeout := flag.Int("t", 2, "Specify the request timeout")
	torRoute := flag.Bool("r", false, "Enable TOR routing")
	dnsLookup := flag.Bool("dns", false, "Use wordlist (-w) and resolve subdomains by querying a DNS")
	dnsLookupCustom := flag.String("dnsC", "", "Specify custom DNS server")
	dnsLookupTimeout := flag.Int("dnsT", 500, "Specify timeout for DNS queries in ms")
	rDnsLookupFilePath := flag.String("rF", "", "IP address list file path")
	httpRequestDelay := flag.Int("rD", 500, "Specify HTTP request delay")
	disableAllOutput := flag.Bool("dO", false, "Disable all output file streams")
	pingSubdomain := flag.Bool("pS", false, "Ping subdomains (privileged execution required)")
	pingCount := flag.Int("pC", 2, "Specify Ping count (default=2)")
	pingFromFile := flag.String("pF", "", "Ping subdomains from file")
	analyseHeaderSingle := flag.Bool("aS", false, "Analyse HTTP header of single subdomain (specified with -s)")
	httpRequestMethod := flag.String("m", "GET", "Method for sending requests (default: GET)")
	showAllHeaders := flag.Bool("aH", false, "Display all headers of HTTP response")
	detectpurpose := flag.Bool("dP", false, "Detect subdomain purpose (Mail, API...)")
	testMisconf := flag.Bool("mT", false, "Test for common weaknesses")
	allowRedirects := flag.Bool("aR", false, "Allow redirects")
	vhostEnum := flag.Bool("vhost", false, "Enable VHost enumeration")
	flag.Parse()
	if flag.NFlag() == 0 || *help {
		fmt.Println(cli.HelpBanner)
		os.Exit(0)
	}
	if *excHttpCodes != "" && !*httpCode || *filtHttpCodes != "" && !*httpCode {
		return shared.Args{}, errors.New("HTTP code filter enabled, but status codes not requested")
	}
	if *domain != "" && !pkg.IsValidDomain(*domain) {
		return shared.Args{}, errors.New("domain verification failed: " + *domain)
	}
	if !*dnsLookup && *dnsLookupCustom != "" {
		return shared.Args{}, errors.New("custom DNS address can only be set when -dns is specified")
	}
	if *dnsLookup && *wordlistPath == "" {
		return shared.Args{}, errors.New("no wordlist specified, dns method cannot be used")
	}
	args := shared.Args{
		Verbose:             *verbose,
		Domain:              *domain,
		Subdomain:           *subdomain,
		NewOutputDirPath:    *newOutputPath,
		HttpCode:            *httpCode,
		WordlistPath:        *wordlistPath,
		ExcHttpCodes:        *excHttpCodes,
		FilHttpCodes:        *filtHttpCodes,
		AnalyzeHeader:       *analyzeHeader,
		PortScan:            *portScan,
		DbExtendPath:        *dbExtendPath,
		Timeout:             *timeout,
		TorRoute:            *torRoute,
		DnsLookup:           *dnsLookup,
		DnsLookupCustom:     *dnsLookupCustom,
		DnsLookupTimeout:    *dnsLookupTimeout,
		HttpRequestDelay:    *httpRequestDelay,
		RDnsLookupFilePath:  *rDnsLookupFilePath,
		DisableAllOutput:    *disableAllOutput,
		PingSubdomain:       *pingSubdomain,
		PingCount:           *pingCount,
		PingSubdomainsFile:  *pingFromFile,
		AnalyseHeaderSingle: *analyseHeaderSingle,
		HttpRequestMethod:   *httpRequestMethod,
		ShowAllHeaders:      *showAllHeaders,
		DetectPurpose:       *detectpurpose,
		MisconfTest:         *testMisconf,
		AllowRedirects:      *allowRedirects,
		IpAddress:           *ipAddresses,
		EnableVHostEnum:     *vhostEnum,
	}
	return args, nil
}
