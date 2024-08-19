package lib

import (
	"Sentinel/lib/utils"
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	fmt.Fprintln(utils.GStdout)
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
	startTime := time.Now()
	obtainedCounter := 0
	allCounter := 0
	var streams utils.FileStreams
	// Ensure that the wordlist specified by the -w flag exists.
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find wordlist: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	lineCount, err := utils.FileCountLines(args.WordlistPath)
	if err != nil {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Failed to count lines of " + args.WordlistPath,
			ExitError:   err,
		})
	}
	wordlistStream, err := os.Open(args.WordlistPath)
	if err != nil {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	defer wordlistStream.Close()
	scanner := bufio.NewScanner(wordlistStream)
	fmt.Println()
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	err = streams.OpenOutputFileStreams(filePaths)
	if err != nil {
		utils.Glogger.Println(err)
	}
	defer streams.CloseOutputFileStreams()
	for scanner.Scan() {
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
			fmt.Println()
			OutputHandler(&streams, client, args, params)
			obtainedCounter++
		}
		allCounter++
		fmt.Fprintf(utils.GStdout, "\rProgress::[%d/%d]", allCounter, lineCount)
		utils.GStdout.Flush()
	}
	if err := scanner.Err(); err != nil {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Scanner failed",
			ExitError:   err,
		})
	}
	utils.Evaluation(startTime, obtainedCounter)
}
