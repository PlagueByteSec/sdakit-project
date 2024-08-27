package main

import (
	"Sentinel/lib"
	"Sentinel/lib/cli"
	"Sentinel/lib/requests"
	"Sentinel/lib/shared"
	"Sentinel/lib/streams"
	"Sentinel/lib/utils"
	"errors"
	"fmt"
)

func main() {
	args, err := cli.CliParser()
	if err != nil {
		utils.SentinelPanic(err)
	}
	if args.Verbose {
		shared.GVerbose = true
	}
	var (
		filePaths *shared.FilePaths = nil
		isExec    int
	)
	utils.InterruptListenerInit()
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
	utils.SentinelPrintBanner(httpClient)
	shared.GDisplayCount = 0
	if args.DisableAllOutput {
		shared.GDisableAllOutput = true
	} else if args.Domain != "" {
		/*
			Initialize the output file paths and create the output
			directory if it does not already exist.
		*/
		filePaths, err = utils.FilePathInit(&args)
		if err != nil {
			utils.SentinelPanic(err)
		}
	}
	fmt.Fprint(shared.GStdout, "[*] Method: ")
	methods := lib.MethodManagerInit()
	for key, method := range methods {
		switch key {
		case shared.Passive:
			if utils.IsPassiveEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				fmt.Fprintln(shared.GStdout)
				method.Action(&args, httpClient, filePaths)
				isExec++
			}
		case shared.Active:
			if utils.IsActiveEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				method.Action(&args, httpClient, filePaths)
				isExec++
			}
		case shared.Dns:
			if utils.IsDnsEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, method.MethodKey)
				method.Action(&args, httpClient, filePaths)
				isExec++
			}
		}
	}
	extern := lib.ExternsManagerInit()
	for key, method := range extern {
		switch key {
		case shared.RDns:
			if utils.IsRDnsEnumeration(&args) {
				fmt.Fprintln(shared.GStdout, shared.RDns)
				method.Action(&args)
				isExec++
			}
		case shared.Ping:
			if utils.IsPingFromFile(&args) {
				fmt.Fprintln(shared.GStdout, shared.Ping)
				method.Action(&args)
				isExec++
			}
		}
	}
	if isExec == 0 {
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
