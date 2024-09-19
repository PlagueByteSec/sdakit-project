package streams

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	utils "github.com/PlagueByteSec/sdakit-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/coreutils/analysis"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/requests"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
	"github.com/PlagueByteSec/sdakit-project/v2/pkg"
)

func WriteJSON(jsonFileName string) error {
	/*
		Write the summary in JSON format to a file. The default
		directory (output) is used if no custom path is specified with the -j flag.
	*/
	bytes, err := json.MarshalIndent(shared.GJsonResult.Subdomains, "", "	")
	if err != nil {
		logging.GLogger.Log(err.Error())
		return err
	}
	if err := os.WriteFile(jsonFileName, bytes, shared.DefaultPermission); err != nil {
		logging.GLogger.Log(err.Error())
		return errors.New("failed to write JSON to: " + jsonFileName)
	}
	return nil
}

func builderAddContent(consoleOutput <-chan string, sb *strings.Builder, wg *sync.WaitGroup) {
	defer wg.Done()
	for output := range consoleOutput {
		sb.WriteString(output)
	}
}

func processFilter(filterValues string) []string {
	delim := ","
	var filter []string
	if !strings.Contains(filterValues, delim) {
		filter = []string{filterValues}
	} else {
		filter = strings.Split(filterValues, delim)
	}
	return filter
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
					logging.GLogger.Log(err.Error())
				}
			}
		}
		// add ipv4 to json structure
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
					logging.GLogger.Log(err.Error())
				}
			}
		}
		// add ipv6 to json structure
		shared.GSubdomBase.IpAddresses.IPv6 = append(
			shared.GSubdomBase.IpAddresses.IPv6,
			net.ParseIP(ip),
		)
	}
}

func optionsSettingsHandler(settings shared.SettingsHandler, outputChan chan<- string) bool {
	var (
		url            string
		httpStatusCode int
		err            error
		size           int
		hash           string
	)
	switch {
	case settings.URL == "":
		url = fmt.Sprintf("http://%s", settings.Params.Subdomain)
	default:
		url = settings.URL
	}
	if settings.Args.HttpCode {
		vhost := settings.Args.EnableVHostEnum
		switch {
		case vhost:
			var body []byte
			_, httpStatusCode, body, err = requests.RequestHandlerCore(&requests.HttpRequestBase{
				HttpClient:             settings.HttpClient,
				CustomUrl:              url,
				HttpMethod:             settings.Args.HttpRequestMethod,
				Subdomain:              settings.Params.Subdomain,
				ResponseNeedStatusCode: true,
				ResponseNeedBody:       true,
			})
			if err != nil {
				logging.GLogger.Log(err.Error())
				return false
			}
			if httpStatusCode == -1 {
				return false
			}
			size = len(body)
			hash = fmt.Sprintf("%x", sha256.Sum256(body))
		default:
			_, httpStatusCode, _, _ = requests.RequestHandlerCore(&requests.HttpRequestBase{
				HttpClient:             settings.HttpClient,
				CustomUrl:              url,
				HttpMethod:             settings.Args.HttpRequestMethod,
				ResponseNeedStatusCode: true,
			})
		}
		statusCodeConv := strconv.Itoa(httpStatusCode)
		switch {
		case httpStatusCode == -1:
			statusCodeConv = shared.NotAvailable
		default:
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
		if vhost {
			trimFilter := processFilter(settings.Args.FilterHttpSize)
			if len(trimFilter) != 0 && pkg.IsInSlice(fmt.Sprintf("%d", size), trimFilter) {
				return false
			}
			outputChan <- fmt.Sprintf(" | Size: %d\n", size)
			outputChan <- fmt.Sprintf(" | Hash: %s\n", hash)
		}
		outputChan <- fmt.Sprintf(" | HTTP Status Code: %s\n", statusCodeConv)
	} else if !settings.Args.DisableAllOutput {
		OutputWrapper(settings.IpAddrs, settings.Params, settings.Streams)
	}
	if settings.Args.AnalyzeHeader {
		headers := requests.AnalyseHttpHeader(settings.HttpClient, settings.Params.Subdomain, settings.Args.HttpRequestMethod)
		outputChan <- headers
	}
	if settings.IpAddrsOut != "" {
		outputChan <- fmt.Sprintf(" | IP Addresses: %s\n", settings.IpAddrsOut)
	}
	if settings.Args.PortScan != "" {
		utils.PortScanWrapper(outputChan, settings.Params.Subdomain, settings.Args.PortScan)
	}
	if settings.Args.PingSubdomain {
		utils.PingWrapper(outputChan, settings.Params.Subdomain, settings.Args.PingCount)
	}
	requests.SetDnsEnumType() // Handle type by global switch
	if settings.Args.DetectPurpose {
		if requests.HttpCodeCheck(settings, url) {
			shared.GShowAllHeaders = true
			headers := requests.AnalyseHttpHeader(settings.HttpClient, settings.Params.Subdomain, settings.Args.HttpRequestMethod)
			check := analysis.SubdomainCheck{
				Subdomain:     settings.Params.Subdomain,
				ConsoleOutput: outputChan,
				HttpHeaders:   headers,
				HttpClient:    settings.HttpClient,
			}
			check.TargetAnalyseHTTP()
		} else {
			// HTTP request failed: run non HTTP tests
			check := analysis.SubdomainCheck{
				Subdomain:     settings.Params.Subdomain,
				ConsoleOutput: outputChan,
			}
			check.TargetAnalyseNonHTTP()
		}
	}
	// httpCodeCheck: do not perform analysis if the HTTP request fails (-1)
	if settings.Args.MisconfTest && requests.HttpCodeCheck(settings, url) {
		check := analysis.SubdomainCheck{
			Subdomain:     settings.Params.Subdomain,
			ConsoleOutput: outputChan,
			HttpClient:    settings.HttpClient,
		}
		check.TestSecurity()
	}
	return true
}

func OutputHandler(streams *shared.FileStreams, client *http.Client, args *shared.Args,
	params shared.Params, url string) {
	if shared.GScanMethod != shared.Passive {
		shared.PoolAppendValue(params.Subdomain, &shared.GPoolBase.PoolSubdomains)
	}
	shared.GStdout.Flush()
	if args.HttpCode || args.AnalyzeHeader {
		time.Sleep(time.Duration(args.HttpRequestDelay) * time.Millisecond)
	}
	/*
		Perform a DNS lookup to determine the IP addresses (IPv4 and IPv6). The addresses will
		be returned as a slice and separated as strings.
	*/
	var (
		ipAddrsOut string
		ipAddrs    []string
	)
	if !utils.IsPassiveEnumeration(args) {
		ipAddrsOut, ipAddrs = utils.IpResolveWrapper(shared.GDnsResolver, params.Subdomain)
		if ipAddrs == nil {
			return
		}
	}
	var wg sync.WaitGroup
	var consoleOutput strings.Builder
	outputChan := make(chan string, 20)
	wg.Add(1)
	go builderAddContent(outputChan, &consoleOutput, &wg)
	if !shared.GDisableAllOutput {
		shared.GSubdomBase = shared.SubdomainBase{}
		shared.GSubdomBase.Subdomain = append(shared.GSubdomBase.Subdomain, params.Subdomain)
	}
	outputChan <- fmt.Sprintf("\r[+] %-100s\n", params.Subdomain)
	/*
		Split the arguments specified by the -f and -e flags by comma.
		The values within the slices will be used to filter the results.
	*/
	codeFilter := processFilter(args.FilHttpCodes)
	codeFilterExc := processFilter(args.ExcHttpCodes)
	pkg.ResetSlice(&codeFilter)
	pkg.ResetSlice(&codeFilterExc)
	succeed := optionsSettingsHandler(shared.SettingsHandler{
		Streams:       streams,
		Args:          args,
		Params:        params,
		HttpClient:    client,
		ConsoleOutput: outputChan,
		CodeFilterExc: codeFilterExc,
		CodeFilter:    codeFilter,
		IpAddrs:       ipAddrs,
		IpAddrsOut:    ipAddrsOut,
		URL:           url,
	}, outputChan)
	if !succeed {
		close(outputChan)
		return
	}
	if !args.DisableAllOutput {
		shared.GJsonResult.Subdomains = append(shared.GJsonResult.Subdomains, shared.GSubdomBase)
	}
	close(outputChan)
	wg.Wait()
	// Display the final result block
	fmt.Fprintln(shared.GStdout, consoleOutput.String())
	shared.GDisplayCount++
}
