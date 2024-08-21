package main

import (
	"Sentinel/lib"
	"Sentinel/lib/utils"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var (
		httpClient   *http.Client
		err          error
		localVersion string
		repoVersion  string
		sigChan      chan os.Signal
		filePaths    *utils.FilePaths = nil
	)
	args, err := lib.CliParser()
	if err != nil {
		goto exitErr
	}
	if args.Verbose {
		utils.GVerbose = true
	}
	/*
		Create a channel to receive interrupt signals from the OS.
		The goroutine continuously listens for an interrupt signal
		(Ctrl+C) and handles the interruption.
	*/
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			utils.SentinelExit(utils.SentinelExitParams{
				ExitCode:    0,
				ExitMessage: "\n\nG0oDBy3!",
				ExitError:   nil,
			})
		}
	}()
	/*
		Set up the HTTP client with a default timeout of 5 seconds
		or a custom timeout specified with the -t flag. If the -r flag
		is provided, the standard HTTP client will be ignored, and
		the client will be configured to route through TOR.
	*/
	httpClient, err = lib.HttpClientInit(&args)
	if err != nil {
		goto exitErr
	}
	localVersion = utils.GetCurrentLocalVersion()
	repoVersion = lib.GetCurrentRepoVersion(httpClient)
	fmt.Fprintf(utils.GStdout, " ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	utils.GStdout.Flush()
	utils.VersionCompare(repoVersion, localVersion)
	utils.GDisplayCount = 0
	/*
		Initialize the output file paths and create the output
		directory if it does not already exist.
	*/
	filePaths, err = utils.FilePathInit(&args)
	if err != nil {
		goto exitErr
	}
	fmt.Fprint(utils.GStdout, "[*] Method: ")
	if args.WordlistPath == "" && args.RDnsLookupFilePath == "" {
		// Perform enumeration using external resources
		fmt.Fprintln(utils.GStdout, "PASSIVE")
		lib.PassiveEnum(&args, httpClient, filePaths)
	}
	if args.DnsLookup && args.WordlistPath != "" {
		// Perform enumeration using DNS
		fmt.Fprintln(utils.GStdout, "DNS")
		lib.DnsEnum(&args, httpClient, filePaths)
	}
	if args.WordlistPath != "" && !args.DnsLookup && args.RDnsLookupFilePath == "" {
		// Perform enumeration using brute force
		fmt.Fprintln(utils.GStdout, "ACTIVE")
		lib.ActiveEnum(&args, httpClient, filePaths)
	}
	if args.RDnsLookupFilePath != "" {
		fmt.Fprintln(utils.GStdout, "RDNS")
		lib.RDnsEnum(&args)
	} else {
		lib.WriteJSON(filePaths.FilePathJSON)
	}
	/*
		Save the summary (including IPv4, IPv6, ports if requested,
		and subdomains) in JSON format within the output directory.
	*/
	utils.SentinelExit(utils.SentinelExitParams{
		ExitCode:    0,
		ExitMessage: "",
		ExitError:   nil,
	})
exitErr:
	fmt.Fprintln(utils.GStdout, err)
	utils.GStdout.Flush()
	utils.Glogger.Println(err)
	utils.Glogger.Fatalf("Program execution failed")
}
