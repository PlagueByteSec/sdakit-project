package lib

import (
	"Sentinel/lib/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func AnalyzeHeaderWrapper(consoleOutput *strings.Builder, ipAddrsOut string,
	client *http.Client, params utils.Params) {
	/*
		Analyze the HTTP header and add the results to the consoleOutput
		string builder if it exists.
	*/
	headers, count := AnalyseHttpHeader(client, params.Subdomain)
	if ipAddrsOut != "" {
		if count != 0 {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╠═[ %s", ipAddrsOut))
		} else {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╚═[ %s\n", ipAddrsOut))
		}
	}
	if headers != "" && count != 0 {
		consoleOutput.WriteString(fmt.Sprintf("\n\t%s\n", headers))
	}
}

func PortScanWrapper(consoleOutput *strings.Builder, params utils.Params, args *utils.Args) {
	ports, err := ScanPortsSubdomain(params.Subdomain, args.PortScan)
	if err != nil {
		utils.Glogger.Println(err)
	}
	if ports != "" {
		consoleOutput.WriteString(ports)
	}
}

func IpResolveWrapper(args *utils.Args, params utils.Params) (string, []string) {
	ipAddrs := utils.RequestIpAddresses(false, params.Subdomain) // DNS lookup
	if ipAddrs == nil {
		// Skip results that cannot be resolved to an IP address
		return "", nil
	}
	ipAddrsOut := strings.Join(ipAddrs, ", ")
	return ipAddrsOut, ipAddrs
}

func WordlistInit(args *utils.Args) (*os.File, int) {
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "could not find wordlist: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	lineCount, err := utils.FileCountLines(args.WordlistPath) // dup
	if err != nil {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Failed to count lines of " + args.WordlistPath,
			ExitError:   err,
		})
	} //
	wordlistStream, err := os.Open(args.WordlistPath) // dup
	if err != nil {
		utils.Glogger.Println(err)
		utils.SentinelExit(utils.SentinelExitParams{
			ExitCode:    -1,
			ExitMessage: "Unable to open stream (read-mode) to: " + args.WordlistPath,
			ExitError:   err,
		})
	}
	return wordlistStream, lineCount
}

func OpenOutputFileStreamsWrapper(filePaths *utils.FilePaths) {
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if err := utils.GStreams.OpenOutputFileStreams(filePaths); err != nil {
		utils.Glogger.Println(err)
	}
}
