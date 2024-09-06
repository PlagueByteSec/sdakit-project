package streams

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/fhAnso/Sentinel/v1/internal/coreutils"
	"github.com/fhAnso/Sentinel/v1/internal/coreutils/analysis"
	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/pkg"
)

func WriteJSON(jsonFileName string) error {
	/*
		Write the summary in JSON format to a file. The default
		directory (output) is used if no custom path is specified with the -j flag.
	*/
	bytes, err := json.MarshalIndent(shared.GJsonResult.Subdomains, "", "	")
	if err != nil {
		shared.Glogger.Println(err)
		return err
	}
	if err := os.WriteFile(jsonFileName, bytes, shared.DefaultPermission); err != nil {
		shared.Glogger.Println(err)
		return errors.New("failed to write JSON to: " + jsonFileName)
	}
	return nil
}

func IpManage(params shared.Params, ip string, fileStream *shared.FileStreams) {
	/*
		Request the IP version based on the given IP address string. A check
		is performed to verify that the address written to the output file is not
		duplicated. If successful, the address will be written to the appropriate output file.
	*/
	ipAddrVersion := pkg.GetIpVersion(ip)
	switch ipAddrVersion {
	case 4:
		params.FileContentIPv4Addrs = ip
		if !pkg.IsInSlice(params.FileContentIPv4Addrs, shared.GPoolBase.PoolIPv4Addresses) {
			shared.GPoolBase.PoolIPv4Addresses = append(shared.GPoolBase.PoolIPv4Addresses, params.FileContentIPv4Addrs)
			if !shared.GDisableAllOutput {
				err := WriteOutputFileStream(fileStream.Ipv4AddrStream, params.FileContentIPv4Addrs)
				if err != nil {
					fileStream.Ipv4AddrStream.Close()
					shared.Glogger.Println(err)
				}
			}
		}
		shared.GSubdomBase.IpAddresses.IPv4 = append(
			shared.GSubdomBase.IpAddresses.IPv4,
			net.ParseIP(ip),
		)
	case 6:
		params.FileContentIPv6Addrs = ip
		if !pkg.IsInSlice(params.FileContentIPv6Addrs, shared.GPoolBase.PoolIPv6Addresses) {
			shared.GPoolBase.PoolIPv6Addresses = append(shared.GPoolBase.PoolIPv6Addresses, params.FileContentIPv6Addrs)
			if !shared.GDisableAllOutput {
				err := WriteOutputFileStream(fileStream.Ipv6AddrStream, params.FileContentIPv6Addrs)
				if err != nil {
					fileStream.Ipv6AddrStream.Close()
					shared.Glogger.Println(err)
				}
			}
		}
		shared.GSubdomBase.IpAddresses.IPv6 = append(
			shared.GSubdomBase.IpAddresses.IPv6,
			net.ParseIP(ip),
		)
	}
}

func optionsSettingsHandler(settings shared.SettingsHandler) bool {
	url := fmt.Sprintf("http://%s", settings.Params.Subdomain)
	if settings.Args.HttpCode {
		httpStatusCode := requests.HttpStatusCode(settings.HttpClient, url, settings.Args.HttpRequestMethod)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if httpStatusCode == -1 {
			statusCodeConv = shared.NotAvailable
		} else {
			shared.PoolAppendValue(settings.Params.Subdomain, &shared.GPoolBase.PoolHttpSuccessSubdomains)
		}
		/*
			Ensure that the status codes are correctly filtered by comparing the
			results with codeFilter and CodeFilterExc.
		*/
		if len(settings.CodeFilter) >= 1 && !pkg.IsInSlice(statusCodeConv, settings.CodeFilter) ||
			len(settings.CodeFilterExc) >= 1 && pkg.IsInSlice(statusCodeConv, settings.CodeFilterExc) {
			return false
		} else if !settings.Args.DisableAllOutput {
			OutputWrapper(settings.IpAddrs, settings.Params, settings.Streams)
		}
		settings.ConsoleOutput.WriteString(fmt.Sprintf(" | HTTP Status Code: %s\n", statusCodeConv))
	} else if !settings.Args.DisableAllOutput {
		OutputWrapper(settings.IpAddrs, settings.Params, settings.Streams)
	}
	if settings.Args.AnalyzeHeader {
		headers := requests.AnalyseHttpHeader(settings.HttpClient, settings.Params.Subdomain, settings.Args.HttpRequestMethod)
		settings.ConsoleOutput.WriteString(headers)
	}
	if settings.IpAddrsOut != "" {
		settings.ConsoleOutput.WriteString(fmt.Sprintf(" | IP Addresses: %s\n", settings.IpAddrsOut))
	}
	if settings.Args.PortScan != "" {
		utils.PortScanWrapper(settings.ConsoleOutput, settings.Params.Subdomain, settings.Args.PortScan)
	}
	if settings.Args.PingSubdomain {
		utils.PingWrapper(settings.ConsoleOutput, settings.Params.Subdomain, settings.Args.PingCount)
	}
	requests.SetDnsEnumType() // Handle type by global switch
	if settings.Args.DetectPurpose {
		settings.ConsoleOutput.WriteString(" | Trying to identify the subdomain purpose...\n")
		shared.GShowAllHeaders = true
		headers := requests.AnalyseHttpHeader(settings.HttpClient, settings.Params.Subdomain, settings.Args.HttpRequestMethod)
		check := analysis.SubdomainCheck{
			Subdomain:     settings.Params.Subdomain,
			ConsoleOutput: settings.ConsoleOutput,
			HttpHeaders:   headers,
			HttpClient:    settings.HttpClient,
		}
		check.Purpose() // run.go
	}
	// httpCodeCheck: do not perform analysis if the HTTP request fails (-1)
	if settings.Args.MisconfTest && requests.HttpCodeCheck(settings, url) {
		settings.ConsoleOutput.WriteString(" | Testing for common weaknesses...\n")
		check := analysis.SubdomainCheck{
			Subdomain:     settings.Params.Subdomain,
			ConsoleOutput: settings.ConsoleOutput,
			HttpClient:    settings.HttpClient,
		}
		check.Misconfigurations() // run.go
	}
	return true
}

func OutputHandler(streams *shared.FileStreams, client *http.Client, args *shared.Args, params shared.Params) {
	if args.HttpCode || args.AnalyzeHeader {
		time.Sleep(time.Duration(args.HttpRequestDelay) * time.Millisecond)
	}
	shared.GStdout.Flush()
	/*
		Perform a DNS lookup to determine the IP addresses (IPv4 and IPv6). The addresses will
		be returned as a slice and separated as strings.
	*/
	ipAddrsOut, ipAddrs := utils.IpResolveWrapper(shared.GDnsResolver, params.Subdomain)
	if ipAddrs == nil {
		return
	}
	var (
		// Use strings.Builder for better output control
		consoleOutput strings.Builder
		codeFilterExc []string
		codeFilter    []string
	)
	if !shared.GDisableAllOutput {
		shared.GSubdomBase = shared.SubdomainBase{}
		shared.GSubdomBase.Subdomain = append(shared.GSubdomBase.Subdomain, params.Subdomain)
	}
	consoleOutput.WriteString(fmt.Sprintf("[+] %-40s\n", params.Subdomain))
	/*
		Split the arguments specified by the -f and -e flags by comma.
		The values within the slices will be used to filter the results.
	*/
	delim := ","
	if !strings.Contains(args.ExcHttpCodes, delim) {
		codeFilterExc = []string{args.ExcHttpCodes}
	} else {
		codeFilterExc = strings.Split(args.ExcHttpCodes, delim)
	}
	if !strings.Contains(args.FilHttpCodes, delim) {
		codeFilter = []string{args.FilHttpCodes}
	} else {
		codeFilter = strings.Split(args.FilHttpCodes, delim)
	}
	pkg.ResetSlice(&codeFilter)
	pkg.ResetSlice(&codeFilterExc)
	succeed := optionsSettingsHandler(shared.SettingsHandler{
		Streams:       streams,
		Args:          args,
		Params:        params,
		HttpClient:    client,
		ConsoleOutput: &consoleOutput,
		CodeFilterExc: codeFilterExc,
		CodeFilter:    codeFilter,
		IpAddrs:       ipAddrs,
		IpAddrsOut:    ipAddrsOut,
	})
	if !succeed {
		return
	}
	if !args.DisableAllOutput {
		shared.GJsonResult.Subdomains = append(shared.GJsonResult.Subdomains, shared.GSubdomBase)
	}
	// Display the final result block
	fmt.Fprintln(shared.GStdout, consoleOutput.String())
	shared.GDisplayCount++
}
