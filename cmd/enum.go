package cmd

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils/analysis"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/requests"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/streams"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

func OpenStreamsEnum(args *shared.Args, filePaths *shared.FilePaths) (*os.File, int) {
	wordlistStream, entryCount := streams.WordlistStreamInit(args)
	if !shared.GDisableAllOutput {
		streams.OpenOutputFileStreamsWrapper(filePaths)
	}
	return wordlistStream, entryCount
}

func NextEntry() {
	shared.GStdout.Flush()
	shared.GObtainedCounter++
}

func Finish(scanner *bufio.Scanner) {
	streams.ScannerCheckError(scanner)
	fmt.Print("\r")
	utils.PrintSummary(shared.GStartTime, shared.GObtainedCounter)
}

func PassiveEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
	shared.GStdout.Flush()
	startTime := time.Now()
	utils.PrintVerbose("[*] Formatting db entries..\n")
	/*
		Read and format the entries listed in db.go, and if specified,
		also handle the endpoints indicated by the -x flag.
	*/
	endpoints, err := utils.EditDbEntries(args)
	if err != nil {
		logging.GLogger.Log(err.Error())
	}
	utils.PrintVerbose("[*] Sending GET request to endpoints..\n")
	/*
		Send a GET request to each endpoint and filter the results. The results will
		be temporarily stored in the appropriate pool. Duplicates will be removed.
	*/
	for idx := 0; idx < len(endpoints); idx++ {
		if err := requests.EndpointRequest(args.HttpRequestMethod, args.Domain, endpoints[idx], client); err != nil {
			logging.GLogger.Log(err.Error())
		}
	}
	poolSize := len(shared.GPoolBase.PoolSubdomains)
	if poolSize == 0 {
		fmt.Fprintln(shared.GStdout, "[-] Could not determine subdomains :(")
		shared.GStdout.Flush()
		os.Exit(0)
	}
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if !args.DisableAllOutput {
		streams.OpenOutputFileStreamsWrapper(filePaths)
		defer streams.CloseOutputFileStreams(&shared.GStreams)
	}
	/*
		Iterate through the subdomain pool and process the current entry. The OutputHandler
		function will ensure that all fetched data is separated and stored within the output
		files, and it will also handle other actions specified by the command line.
	*/
	utils.PrintVerbose("[*] Retrieved %d subdomains: Evaluation.\n", poolSize)
	for _, subdomain := range shared.GPoolBase.PoolSubdomains {
		paramsSetupFiles := shared.ParamsSetupFilesBase{
			FileParams: &shared.Params{},
			CliArgs:    args,
			FilePaths:  filePaths,
			Subdomain:  subdomain,
		}
		streams.ParamsSetupFiles(paramsSetupFiles)
		streams.OutputHandlerWrapper(subdomain, client, args, &paramsSetupFiles, "")
	}
	// Evaluate the summary and format it for writing to stdout.
	utils.PrintSummary(startTime, poolSize)
}

func DirectEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
	wordlistStream, entryCount := streams.WordlistStreamInit(args)
	defer wordlistStream.Close()
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Fprintln(shared.GStdout)
	if !shared.GDisableAllOutput {
		streams.OpenOutputFileStreamsWrapper(filePaths)
		defer streams.CloseOutputFileStreams(&shared.GStreams)
	}
	for scanner.Scan() {
		shared.GSubdomBase = shared.SubdomainBase{}
		entry := scanner.Text()
		if pkg.LineIgnore(entry) {
			continue
		}
		url := fmt.Sprintf("http://%s.%s", entry, args.Domain)
		_, statusCode, _, _ := requests.RequestHandlerCore(&requests.HttpRequestBase{
			HttpClient:             client,
			CustomUrl:              url,
			HttpMethod:             args.HttpRequestMethod,
			ResponseNeedStatusCode: true,
		})
		/*
			Skip failed GET requests and set the successful response subdomains to the
			Params struct. The OutputHandler function will ensure that all fetched data
			is separated and stored within the output files, and it will also handle
			other actions specified by the command line.
		*/
		if statusCode != -1 {
			subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
			paramsSetupFiles := shared.ParamsSetupFilesBase{
				FileParams: &shared.Params{},
				CliArgs:    args,
				FilePaths:  filePaths,
				Subdomain:  subdomain,
			}
			streams.ParamsSetupFiles(paramsSetupFiles)
			fmt.Fprint(shared.GStdout, "\r")
			streams.OutputHandlerWrapper(subdomain, client, args, &paramsSetupFiles, "")
			NextEntry()
		}
		utils.PrintProgress(entryCount)
	}
	Finish(scanner)
}

func DnsEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
	/*
		Ensure that the specified wordlist can be found and open
		a read-only stream.
	*/
	wordlistStream, entryCount := OpenStreamsEnum(args, filePaths)
	defer wordlistStream.Close()
	defer streams.CloseOutputFileStreams(&shared.GStreams)
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Fprintln(shared.GStdout)
	/*
		Check if a custom DNS server address is specified by the -dnsC
		flag. If it is specified, ensure that the IP address follows the
		correct pattern and that the specified port is within the correct range.
	*/
	if args.DnsLookupCustom != "" {
		testValue := strings.Split(args.DnsLookupCustom, ":")
		dnsServerIp := net.ParseIP(testValue[0])
		if testValue == nil {
			utils.SentinelExit(shared.SentinelExitParams{
				ExitCode:    -1,
				ExitMessage: "Please specify a valid DNS server address",
				ExitError:   nil,
			})
		}
		dnsServerPort, err := strconv.ParseInt(testValue[1], 0, 16)
		if err != nil || dnsServerPort < 1 && dnsServerPort > 65535 {
			utils.SentinelExit(shared.SentinelExitParams{
				ExitCode:    -1,
				ExitMessage: "Please specify a valid DNS server port",
				ExitError:   nil,
			})
		}
		shared.CustomDnsServer = string(dnsServerIp)
	}
	for scanner.Scan() {
		shared.GDnsResults = []string{}
		entry := scanner.Text()
		if pkg.LineIgnore(entry) {
			continue
		}
		subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
		requests.SetDnsEnumType()
		/*
			Perform DNS lookup against the current subdomain. The results
			will be stored in GDnsResults.
		*/
		requests.DnsLookups(shared.GDnsResolver, shared.DnsLookupOptions{
			IpAddress: nil,
			Subdomain: subdomain,
		})
		if len(shared.GDnsResults) != 0 {
			paramsSetupFiles := shared.ParamsSetupFilesBase{
				FileParams: &shared.Params{},
				CliArgs:    args,
				FilePaths:  filePaths,
				Subdomain:  subdomain,
			}
			streams.ParamsSetupFiles(paramsSetupFiles)
			streams.OutputHandlerWrapper(subdomain, client, args, &paramsSetupFiles, "")
			NextEntry()
		}
		utils.PrintProgress(entryCount)
	}
	Finish(scanner)
}

func VHostEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
	shared.GStdout.Flush()
	wordlistStream, entryCount := OpenStreamsEnum(args, filePaths)
	defer wordlistStream.Close()
	defer streams.CloseOutputFileStreams(&shared.GStreams)
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Fprintln(shared.GStdout)
	ipAddress := net.ParseIP(args.IpAddress).String()
	portScanSummary, _, err := requests.ScanPortRange(ipAddress, "80,443", false)
	if err != nil {
		fmt.Fprintln(shared.GStdout, err.Error())
		return
	}
	const (
		HTTP             = "HTTP"
		AlternativeHTTP  = "AlternativeHTTP"
		HTTPS            = "HTTPS"
		AlternativeHTTPS = "AlternativeHTTPS"
	)
	port := map[string]string{
		HTTP:             "80",
		AlternativeHTTP:  "8080",
		HTTPS:            "443",
		AlternativeHTTPS: "8443",
	}
	var proto analysis.HTTP
	if strings.Contains(portScanSummary, port[HTTPS]) || strings.Contains(portScanSummary, port[AlternativeHTTPS]) {
		proto = analysis.HTTP(analysis.Basic)
	} else if strings.Contains(portScanSummary, port[HTTP]) || strings.Contains(portScanSummary, port[AlternativeHTTP]) {
		proto = analysis.HTTP(analysis.Secure)
	} else {
		fmt.Fprintf(shared.GStdout, "[-] Port scan failed, scanned: ")
		for _, value := range port {
			fmt.Fprintf(shared.GStdout, value+" ")
		}
		fmt.Fprintln(shared.GStdout)
		return
	}
	ipUrl := analysis.MakeUrl(proto, ipAddress)
	_, statusCode, _, _ := requests.RequestHandlerCore(&requests.HttpRequestBase{
		HttpClient:             client,
		CustomUrl:              ipUrl,
		HttpMethod:             args.HttpRequestMethod,
		ResponseNeedStatusCode: true,
	})
	if statusCode == -1 {
		fmt.Fprintf(shared.GStdout, "[-] %s: no response, abort.\n", ipUrl)
		return
	}
	for scanner.Scan() {
		shared.GSubdomBase = shared.SubdomainBase{}
		entry := scanner.Text()
		if pkg.LineIgnore(entry) {
			continue
		}
		utils.PrintProgress(entryCount)
		subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
		_, statusCode, _, _ := requests.RequestHandlerCore(&requests.HttpRequestBase{
			HttpClient:             client,
			CustomUrl:              ipUrl,
			HttpMethod:             args.HttpRequestMethod,
			Subdomain:              subdomain,
			ResponseNeedStatusCode: true,
		})
		if statusCode == -1 {
			continue
		}
		args.HttpCode = true
		paramsSetupFiles := shared.ParamsSetupFilesBase{
			FileParams: &shared.Params{},
			CliArgs:    args,
			FilePaths:  filePaths,
			Subdomain:  subdomain,
		}
		streams.ParamsSetupFiles(paramsSetupFiles)
		streams.OutputHandlerWrapper(subdomain, client, args, &paramsSetupFiles, ipUrl)
	}
	Finish(scanner)
}
