package lib

import (
	"Sentinel/lib/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteJSON(jsonFileName string) error {
	/*
		Write the summary in JSON format to a file. The default
		directory (output) is used if no custom path is specified with the -j flag.
	*/
	bytes, err := json.MarshalIndent(utils.GJsonResult.Subdomains, "", "	")
	if err != nil {
		utils.Glogger.Println(err)
		return err
	}
	if err := os.WriteFile(jsonFileName, bytes, utils.DefaultPermission); err != nil {
		utils.Glogger.Println(err)
		return errors.New("failed to write JSON to: " + jsonFileName)
	}
	return nil
}

func OutputHandler(streams *utils.FileStreams, client *http.Client, args *utils.Args, params utils.Params) {
	if args.HttpCode || args.AnalyzeHeader {
		time.Sleep(time.Duration(args.HttpRequestDelay) * time.Millisecond)
	}
	utils.GStdout.Flush()
	/*
		Perform a DNS lookup to determine the IP addresses (IPv4 and IPv6). The addresses will
		be returned as a slice and separated as strings.
	*/
	ipAddrsOut, ipAddrs := IpResolveWrapper(args, params)
	if ipAddrs == nil {
		return
	}
	var (
		// Use strings.Builder for better output control
		consoleOutput strings.Builder
		err           error
		tempExclude   string
	)
	utils.GSubdomBase.Subdomain = append(utils.GSubdomBase.Subdomain, params.Subdomain)
	for _, ip := range ipAddrs {
		utils.IpManage(params, ip, streams)
	}
	err = utils.WriteOutputFileStream(streams.SubdomainStream, params.FileContentSubdoms)
	if err != nil {
		streams.SubdomainStream.Close()
	}
	consoleOutput.WriteString(fmt.Sprintf(" ══[ %s", params.Subdomain))
	/*
		Split the arguments specified by the -f and -e flags by comma.
		The values within the slices will be used to filter the results.
	*/
	delim := ","
	if !strings.Contains(args.ExcHttpCodes, delim) {
		tempExclude = args.ExcHttpCodes
	}
	codeFilter := strings.Split(args.FilHttpCodes, delim)
	codeFilterExc := strings.Split(args.ExcHttpCodes, delim)
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Subdomain)
		httpStatusCode := HttpStatusCode(client, url)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if httpStatusCode == -1 {
			statusCodeConv = utils.NotAvailable
		}
		/*
			Ensure that the status codes are correctly filtered by comparing the
			results with codeFilter and CodeFilterExc.
		*/
		if len(codeFilter) != 1 && !utils.InArgList(statusCodeConv, codeFilter) {
			return
		}
		if len(codeFilterExc) != 1 && utils.InArgList(statusCodeConv, codeFilterExc) {
			return
		}
		if tempExclude == statusCodeConv {
			return
		}
		consoleOutput.WriteString(fmt.Sprintf(", HTTP Status Code: %s", statusCodeConv))
	}
	if args.AnalyzeHeader {
		AnalyzeHeaderWrapper(&consoleOutput, ipAddrsOut, client, params)
	} else {
		if ipAddrsOut != "" {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╚═[ %s\n", ipAddrsOut))
		}
	}
	if args.PortScan != "" {
		PortScanWrapper(&consoleOutput, params, args)
	}
	utils.GJsonResult.Subdomains = append(utils.GJsonResult.Subdomains, utils.GSubdomBase)
	// Display the final result block
	fmt.Fprintln(utils.GStdout, consoleOutput.String())
	utils.GDisplayCount++
}
