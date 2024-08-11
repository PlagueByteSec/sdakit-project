package lib

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func evaluation(startTime time.Time, count int) {
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("\n[*] %d subdomains obtained, %d displayed\n", count, DisplayCount)
	fmt.Printf("[*] Finished in %s\n", duration)
}

// Pool init and preparation
func PassiveEnum(args *Args, client *http.Client) {
	startTime := time.Now()
	if args.Verbose {
		fmt.Println("[*] Formatting db entries..")
	}
	endpoints, err := EditDbEntries(args)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[*] Sending GET request to endpoints..")
	for idx := 0; idx < len(endpoints); idx++ {
		if err := EndpointRequest(client, args.Host, endpoints[idx]); err != nil {
			fmt.Printf("[-] %s\n", err)
		}
	}
	if len(PoolDomains) == 0 {
		fmt.Println("[-] Could not determine subdomains :(")
		os.Exit(0)
	}
	var (
		filePath     string
		filePathIPv4 string
		filePathIPv6 string
	)
	if args.OutFile == "defaultSd" {
		filePath = filepath.Join("output", DefaultOutputName(args.Host))
	} else {
		filePath = args.OutFile
	}
	if args.OutFileIPv4 == "defaultV4" {
		filePathIPv4 = filepath.Join("output", "IPv4-"+DefaultOutputName(args.Host))
	} else {
		filePathIPv4 = args.OutFileIPv4
	}
	if args.OutFileIPv6 == "defaultV6" {
		filePathIPv6 = filepath.Join("output", "IPv6-"+DefaultOutputName(args.Host))
	} else {
		filePathIPv6 = args.OutFileIPv6
	}
	fmt.Println()
	for _, result := range PoolDomains {
		params := Params{
			FilePath:     filePath,
			FilePathIPv4: filePathIPv4,
			FilePathIPv6: filePathIPv6,
			FileContent:  result,
			Result:       result,
			Hostname:     args.Host,
		}
		OutputHandler(client, args, params)
	}
	poolSize := len(PoolDomains)
	evaluation(startTime, poolSize)
}

func DirectEnum(args *Args, client *http.Client) error {
	var counter int
	startTime := time.Now()
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		return errors.New("could not find wordlist: " + args.WordlistPath)
	}
	stream, err := os.Open(args.WordlistPath)
	if err != nil {
		return errors.New("unable to open file stream to wordlist")
	}
	defer stream.Close()
	codeFilter := strings.Split(args.FilHttpCodes, ",")
	codeFilterExc := strings.Split(args.ExcHttpCodes, ",")
	scanner := bufio.NewScanner(stream)
	fmt.Println()
	for scanner.Scan() {
		entry := scanner.Text()
		url := fmt.Sprintf("http://%s.%s", entry, args.Host)
		statusCode := HttpStatusCode(client, url)
		code := strconv.Itoa(statusCode)
		if len(codeFilter) != 0 && !InArgList(code, codeFilter) {
			continue
		}
		if len(codeFilterExc) != 0 && InArgList(code, codeFilterExc) {
			continue
		}
		fmt.Printf(" ===[ %s.%s: %d\n", entry, args.Host, statusCode)
		counter++
	}
	if err := scanner.Err(); err != nil {
		return errors.New("scanner returns an error while reading wordlist")
	}
	evaluation(startTime, counter)
	return nil
}
