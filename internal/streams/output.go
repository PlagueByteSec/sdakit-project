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

	"github.com/fhAnso/Sentinel/v1/internal/requests"
	"github.com/fhAnso/Sentinel/v1/internal/shared"
	"github.com/fhAnso/Sentinel/v1/internal/utils"
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
	consoleOutput.WriteString(fmt.Sprintf("[+] %s\n", params.Subdomain))
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
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Subdomain)
		httpStatusCode := requests.HttpStatusCode(client, url)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if httpStatusCode == -1 {
			statusCodeConv = shared.NotAvailable
		}
		/*
			Ensure that the status codes are correctly filtered by comparing the
			results with codeFilter and CodeFilterExc.
		*/
		if len(codeFilter) >= 1 && !pkg.IsInSlice(statusCodeConv, codeFilter) ||
			len(codeFilterExc) >= 1 && pkg.IsInSlice(statusCodeConv, codeFilterExc) {
			return
		} else if !args.DisableAllOutput {
			OutputWrapper(ipAddrs, params, streams)
		}
		consoleOutput.WriteString(fmt.Sprintf(" | HTTP Status Code: %s\n", statusCodeConv))
	} else if !args.DisableAllOutput {
		OutputWrapper(ipAddrs, params, streams)
	}
	if args.AnalyzeHeader {
		headers := requests.AnalyseHttpHeader(client, params.Subdomain)
		consoleOutput.WriteString(headers)
	}
	if ipAddrsOut != "" {
		consoleOutput.WriteString(fmt.Sprintf(" | IP Addresses: %s\n", ipAddrsOut))
	}
	if args.PortScan != "" {
		utils.PortScanWrapper(&consoleOutput, params.Subdomain, args.PortScan)
	}
	if args.PingSubdomain {
		utils.PingWrapper(&consoleOutput, params.Subdomain, args.PingCount)
	}
	if !args.DisableAllOutput {
		shared.GJsonResult.Subdomains = append(shared.GJsonResult.Subdomains, shared.GSubdomBase)
	}
	// Display the final result block
	fmt.Fprintln(shared.GStdout, consoleOutput.String())
	shared.GDisplayCount++
}
