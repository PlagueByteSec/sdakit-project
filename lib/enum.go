package lib

import (
	"Sentinel/lib/utils"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func PassiveEnum(args *utils.Args, client *http.Client, filePaths *utils.FilePaths) {
	utils.GStdout.Flush()
	startTime := time.Now()
	utils.VerbosePrint("[*] Formatting db entries..\n")
	/*
		Read and format the entries listed in db.go, and if specified,
		also handle the endpoints indicated by the -x flag.
	*/
	endpoints, err := utils.EditDbEntries(args)
	if err != nil {
		utils.Glogger.Println(err)
	}
	utils.VerbosePrint("[*] Sending GET request to endpoints..\n")
	/*
		Send a GET request to each endpoint and filter the results. The results will
		be temporarily stored in the appropriate pool. Duplicates will be removed.
	*/
	for idx := 0; idx < len(endpoints); idx++ {
		if err := EndpointRequest(client, args.Domain, endpoints[idx]); err != nil {
			utils.Glogger.Println(err)
		}
	}
	if len(utils.GPool.PoolSubdomains) == 0 {
		fmt.Fprintln(utils.GStdout, "[-] Could not determine subdomains :(")
		utils.GStdout.Flush()
		os.Exit(0)
	}
	var streams utils.FileStreams
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	err = streams.OpenOutputFileStreams(filePaths)
	if err != nil {
		utils.Glogger.Println(err)
	}
	defer streams.CloseOutputFileStreams()
	/*
		Iterate through the subdomain pool and process the current entry. The OutputHandler
		function will ensure that all fetched data is separated and stored within the output
		files, and it will also handle other actions specified by the command line.
	*/
	for _, subdomain := range utils.GPool.PoolSubdomains {
		utils.GSubdomBase = utils.SubdomainBase{}
		params := utils.Params{
			Domain:             args.Domain,
			Subdomain:          subdomain,
			FilePathSubdomains: filePaths.FilePathSubdomain,
			FileContentSubdoms: subdomain,
			FilePathIPv4Addrs:  filePaths.FilePathIPv4,
			FilePathIPv6Addrs:  filePaths.FilePathIPv6,
		}
		OutputHandler(&streams, client, args, params)
	}
	poolSize := len(utils.GPool.PoolSubdomains)
	// Evaluate the summary and format it for writing to stdout.
	utils.Evaluation(startTime, poolSize)
}

func ActiveEnum(args *utils.Args, client *http.Client, filePaths *utils.FilePaths) {
	wordlistStream, entryCount := WordlistInit(args)
	defer wordlistStream.Close()
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Fprintln(utils.GStdout)
	OpenOutputFileStreamsWrapper(filePaths)
	defer utils.GStreams.CloseOutputFileStreams()
	for scanner.Scan() {
		utils.GSubdomBase = utils.SubdomainBase{}
		entry := scanner.Text()
		url := fmt.Sprintf("http://%s.%s", entry, args.Domain)
		statusCode := HttpStatusCode(client, url)
		/*
			Skip failed GET requests and set the successful response subdomains to the
			Params struct. The OutputHandler function will ensure that all fetched data
			is separated and stored within the output files, and it will also handle
			other actions specified by the command line.
		*/
		if statusCode != -1 {
			utils.GSubdomBase = utils.SubdomainBase{}
			subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
			params := utils.Params{
				Domain:             args.Domain,
				Subdomain:          subdomain,
				FilePathSubdomains: filePaths.FilePathSubdomain,
				FileContentSubdoms: subdomain,
				FilePathIPv4Addrs:  filePaths.FilePathIPv4,
				FilePathIPv6Addrs:  filePaths.FilePathIPv6,
			}
			fmt.Fprintln(utils.GStdout)
			OutputHandler(&utils.GStreams, client, args, params)
			utils.GStdout.Flush()
			utils.GObtainedCounter++
		}
		utils.PrintProgress(entryCount)
	}
	utils.ScannerCheckError(scanner)
	utils.Evaluation(utils.GStartTime, utils.GObtainedCounter)
}

func DnsEnum(args *utils.Args, client *http.Client, filePaths *utils.FilePaths) {
	/*
		Ensure that the specified wordlist can be found and open
		a read-only stream.
	*/
	wordlistStream, entryCount := WordlistInit(args)
	defer wordlistStream.Close()
	OpenOutputFileStreamsWrapper(filePaths)
	defer utils.GStreams.CloseOutputFileStreams()
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Fprintln(utils.GStdout)
	/*
		Check if a custom DNS server address is specified by the -dnsC
		flag. If it is specified, ensure that the IP address follows the
		correct pattern and that the specified port is within the correct range.
	*/
	if args.DnsLookupCustom != "" {
		testValue := strings.Split(args.DnsLookupCustom, ":")
		dnsServerIp := net.ParseIP(testValue[0])
		if testValue == nil {
			utils.SentinelExit(utils.SentinelExitParams{
				ExitCode:    -1,
				ExitMessage: "Please specify a valid DNS server address",
				ExitError:   nil,
			})
		}
		dnsServerPort, err := strconv.ParseInt(testValue[1], 0, 16)
		if err != nil || dnsServerPort < 1 && dnsServerPort > 65535 {
			utils.SentinelExit(utils.SentinelExitParams{
				ExitCode:    -1,
				ExitMessage: "Please specify a valid DNS server port",
				ExitError:   nil,
			})
		}
		utils.CustomDnsServer = string(dnsServerIp)
	}
	var queryDNS []string
	for scanner.Scan() {
		utils.GSubdomBase = utils.SubdomainBase{}
		entry := scanner.Text()
		subdomain := fmt.Sprintf("%s.%s", entry, args.Domain)
		queryDNS = utils.RequestIpAddresses(false, subdomain)
		if utils.CustomDnsServer != "" {
			// Use custom DNS server address
			queryDNS = utils.RequestIpAddresses(true, subdomain)
		}
		if len(queryDNS) != 0 {
			params := utils.Params{
				Domain:             args.Domain,
				Subdomain:          subdomain,
				FilePathSubdomains: filePaths.FilePathSubdomain,
				FileContentSubdoms: subdomain,
				FilePathIPv4Addrs:  filePaths.FilePathIPv4,
				FilePathIPv6Addrs:  filePaths.FilePathIPv6,
			}
			fmt.Fprintln(utils.GStdout)
			OutputHandler(&utils.GStreams, client, args, params)
			utils.GStdout.Flush()
			utils.GObtainedCounter++
		}
		utils.PrintProgress(entryCount)
	}
	utils.ScannerCheckError(scanner)
	utils.Evaluation(utils.GStartTime, utils.GObtainedCounter)
}
