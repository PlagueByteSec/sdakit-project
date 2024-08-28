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

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/internal/streams"
	"github.com/fhAnso/Sentinel/v1/internal/utils"
)

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
		shared.Glogger.Println(err)
	}
	utils.PrintVerbose("[*] Sending GET request to endpoints..\n")
	/*
		Send a GET request to each endpoint and filter the results. The results will
		be temporarily stored in the appropriate pool. Duplicates will be removed.
	*/
	for idx := 0; idx < len(endpoints); idx++ {
		if err := requests.EndpointRequest(args.HttpRequestMethod, args.Domain, endpoints[idx], client); err != nil {
			shared.Glogger.Println(err)
		}
	}
	if len(shared.GPoolBase.PoolSubdomains) == 0 {
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
	for _, subdomain := range shared.GPoolBase.PoolSubdomains {
		paramsSetupFiles := shared.ParamsSetupFilesBase{
			FileParams: &shared.Params{},
			CliArgs:    args,
			FilePaths:  filePaths,
			Subdomain:  subdomain,
		}
		streams.ParamsSetupFiles(paramsSetupFiles)
		streams.OutputHandler(&shared.GStreams, client, args, *paramsSetupFiles.FileParams)
	}
	poolSize := len(shared.GPoolBase.PoolSubdomains)
	// Evaluate the summary and format it for writing to stdout.
	utils.PrintEvaluation(startTime, poolSize)
}

func ActiveEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
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
		url := fmt.Sprintf("http://%s.%s", entry, args.Domain)
		statusCode := requests.HttpStatusCode(client, url, args.HttpRequestMethod)
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
			streams.OutputHandler(&shared.GStreams, client, args, *paramsSetupFiles.FileParams)
			shared.GStdout.Flush()
			shared.GObtainedCounter++
		}
		utils.PrintProgress(entryCount)
	}
	streams.ScannerCheckError(scanner)
	fmt.Print("\r")
	utils.PrintEvaluation(shared.GStartTime, shared.GObtainedCounter)
}

func DnsEnum(args *shared.Args, client *http.Client, filePaths *shared.FilePaths) {
	/*
		Ensure that the specified wordlist can be found and open
		a read-only stream.
	*/
	wordlistStream, entryCount := streams.WordlistStreamInit(args)
	defer wordlistStream.Close()
	if !shared.GDisableAllOutput {
		streams.OpenOutputFileStreamsWrapper(filePaths)
		defer streams.CloseOutputFileStreams(&shared.GStreams)
	}
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
		subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
		shared.GDnsResolver = requests.DnsResolverInit(false)
		if shared.CustomDnsServer != "" {
			// Use custom DNS server address
			shared.GDnsResolver = requests.DnsResolverInit(true)
		}
		// Perform DNS lookup against the current subdomain
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
			fmt.Fprint(shared.GStdout, "\r")
			streams.OutputHandler(&shared.GStreams, client, args, *paramsSetupFiles.FileParams)
			shared.GStdout.Flush()
			shared.GObtainedCounter++
		}
		utils.PrintProgress(entryCount)
	}
	streams.ScannerCheckError(scanner)
	fmt.Print("\r")
	utils.PrintEvaluation(shared.GStartTime, shared.GObtainedCounter)
}
