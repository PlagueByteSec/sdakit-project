package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fhAnso/Sentinel/v1/internal/cli"
	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/internal/streams"
	"github.com/fhAnso/Sentinel/v1/internal/utils"
)

func methodManager(args shared.Args, httpClient *http.Client, filePaths *shared.FilePaths) {
	// Manager for subdomain enumeration methods that require and HTTP client
	methods := MethodManagerInit()
	for key, method := range methods {
		switch key {
		case shared.Passive: // Request endpoints (certificate transparency logs etc.)
			if utils.IsPassiveEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				fmt.Fprintln(shared.GStdout)
				method.Action(&args, httpClient, filePaths)
				shared.GIsExec++
			}
		case shared.Active: // Brute-force by evaluating HTTP codes
			if utils.IsActiveEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				method.Action(&args, httpClient, filePaths)
				shared.GIsExec++
			}
		case shared.Dns: // Try to resolve a list of subdomains to IP addresses
			if utils.IsDnsEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				method.Action(&args, httpClient, filePaths)
				shared.GIsExec++
			}
		}
	}
	// Manager for commands that require (.txt) lists containing addresses
	extern := ValidsManagerInit()
	for key, method := range extern {
		switch key {
		case shared.RDns: // Resolving addresses from IP list
			if utils.IsRDnsEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, shared.RDns)
				method.Action(&args)
				shared.GIsExec++
			}
		case shared.Ping: // Ping subdomains from subdomain list
			if utils.IsPingFromFile(&args) {
				fmt.Fprintln(shared.GStdout, shared.Ping)
				method.Action(&args)
				shared.GIsExec++
			}
		case shared.HeaderAnalysis:
			if utils.IsHttpHeaderAnalysis(&args) {
				fmt.Fprintln(shared.GStdout, shared.HeaderAnalysis)
				method.Action(&args)
				shared.GIsExec++
			}
		}
	}
}

func Run(args shared.Args) {
	if args.Verbose {
		shared.GVerbose = true
	}
	var filePaths *shared.FilePaths = nil
	InterruptListenerInit()
	/*
		Set up the HTTP client with a default timeout of 5 seconds
		or a custom timeout specified with the -t flag. If the -r flag
		is provided, the standard HTTP client will be ignored, and
		the client will be configured to route through TOR.
	*/
	httpClient, err := requests.HttpClientInit(&args)
	if err != nil {
		utils.SentinelPanic(err)
	}
	// Print banner and compare local with repo version
	utils.PrintBanner(httpClient)
	shared.GDisplayCount = 0
	if args.DisableAllOutput {
		shared.GDisableAllOutput = true
	} else if args.Domain != "" {
		/*
			Initialize the output file paths and create the output
			directory if it does not already exist.
		*/
		filePaths, err = streams.FilePathInit(&args)
		if err != nil {
			utils.SentinelPanic(err)
		}
	}
	fmt.Fprint(shared.GStdout, "[*] Method: ")
	methodManager(args, httpClient, filePaths)
	if shared.GIsExec == 0 {
		utils.SentinelExit(shared.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: cli.HelpBanner,
			ExitError:   errors.New("failed to start enum: syntax error"),
		})
	}
	if !shared.GDisableAllOutput {
		streams.WriteJSON(filePaths.FilePathJSON)
	}
	/*
		Save the summary (including IPv4, IPv6, ports if requested,
		and subdomains) in JSON format within the output directory.
	*/
	utils.SentinelExit(shared.SentinelExitParams{
		ExitCode:    0,
		ExitMessage: "",
		ExitError:   nil,
	})
}
