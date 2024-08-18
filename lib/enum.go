package lib

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func PassiveEnum(args *Args, client *http.Client) {
	GStdout.Flush()
	startTime := time.Now()
	VerbosePrint("[*] Formatting db entries..\n")
	endpoints, err := EditDbEntries(args)
	if err != nil {
		Logger.Println(err)
	}
	VerbosePrint("[*] Sending GET request to endpoints..\n")
	fmt.Fprintln(GStdout)
	for idx := 0; idx < len(endpoints); idx++ {
		if err := EndpointRequest(client, args.Host, endpoints[idx]); err != nil {
			Logger.Println(err)
		}
	}
	if len(GPool.PoolDomains) == 0 {
		fmt.Fprintln(GStdout, "[-] Could not determine subdomains :(")
		os.Exit(0)
	}
	var streams FileStreams
	filePaths := FilePathInit(args)
	err = streams.OpenOutputFileStreams(filePaths)
	if err != nil {
		Logger.Println(err)
	}
	defer streams.CloseOutputFileStreams()
	for _, result := range GPool.PoolDomains {
		params := Params{
			FilePath:     filePaths.FilePathSubdomain,
			FilePathIPv4: filePaths.FilePathIPv4,
			FilePathIPv6: filePaths.FilePathIPv6,
			FileContent:  result,
			Result:       result,
			Hostname:     args.Host,
		}
		OutputHandler(&streams, client, args, params)
	}
	poolSize := len(GPool.PoolDomains)
	Evaluation(startTime, poolSize)
	GPool.PoolCleanup()
}

func DirectEnum(args *Args, client *http.Client) error {
	obtainedCounter := 0
	allCounter := 0
	startTime := time.Now()
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		Logger.Println(err)
		return errors.New("could not find wordlist: " + args.WordlistPath)
	}
	lineCount, err := FileCountLines(args.WordlistPath)
	if err != nil {
		Logger.Println(err)
		return errors.New("Failed to count lines of " + args.WordlistPath)
	}
	stream, err := os.Open(args.WordlistPath)
	if err != nil {
		Logger.Println(err)
		return errors.New("unable to open file stream to wordlist")
	}
	defer stream.Close()
	codeFilter := strings.Split(args.FilHttpCodes, ",")
	codeFilterExc := strings.Split(args.ExcHttpCodes, ",")
	scanner := bufio.NewScanner(stream)
	fmt.Println()
	var streams FileStreams
	filePaths := FilePathInit(args)
	err = streams.OpenOutputFileStreams(filePaths)
	if err != nil {
		Logger.Println(err)
	}
	defer streams.CloseOutputFileStreams()
	for scanner.Scan() {
		entry := scanner.Text()
		url := fmt.Sprintf("http://%s.%s", entry, args.Host)
		statusCode := HttpStatusCode(client, url)
		code := strconv.Itoa(statusCode)
		if statusCode != -1 {
			if len(codeFilter) != 1 && !InArgList(code, codeFilter) {
				continue
			}
			if len(codeFilterExc) != 1 && InArgList(code, codeFilterExc) {
				continue
			}
			subdomain := fmt.Sprintf("%s.%s", entry, args.Host)
			params := Params{
				FilePath:     filePaths.FilePathSubdomain,
				FilePathIPv4: filePaths.FilePathIPv4,
				FilePathIPv6: filePaths.FilePathIPv6,
				FileContent:  subdomain,
				Result:       subdomain,
				Hostname:     args.Host,
			}
			fmt.Println()
			OutputHandler(&streams, client, args, params)
			obtainedCounter++
		}
		allCounter++
		fmt.Fprintf(GStdout, "\rProgress::[%d/%d]", allCounter, lineCount)
		GStdout.Flush()
	}
	if err := scanner.Err(); err != nil {
		Logger.Println(err)
		return errors.New("scanner returns an error while reading wordlist")
	}
	Evaluation(startTime, obtainedCounter)
	GPool.PoolCleanup()
	return nil
}
