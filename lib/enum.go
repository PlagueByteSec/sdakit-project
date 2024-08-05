package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Evaluation(startTime time.Time, count int) {
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("\n[*] %d subdomains obtained, %d displayed\n", count, DisplayCount)
	fmt.Printf("[*] Finished in %s\n", duration)
}

// Pool init and preparation
func PassiveEnum(args *Args) {
	startTime := time.Now()
	pool := make(Pool)
	fmt.Println("[*] Formatting db entries..")
	endpoints := EditDbEntries(args)
	fmt.Println("[*] Sending GET request to endpoints..")
	for idx := 0; idx < len(endpoints); idx++ {
		if err := Request(pool, args.Host, endpoints[idx]); err != nil {
			fmt.Printf("[-] %s\n", err)
			continue
		}
	}
	if len(pool) == 0 {
		fmt.Println("[-] Could not determine subdomains :(")
		os.Exit(0)
	}
	fmt.Println()
	for result := range pool {
		var (
			filePath     string
			filePathIPv4 string
			filePathIPv6 string
		)
		if args.OutFile == "default" {
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
		params := Params{
			FilePath:     filePath,
			FilePathIPv4: filePathIPv4,
			FilePathIPv6: filePathIPv6,
			FileContent:  result,
			Result:       result,
			Hostname:     args.Host,
		}
		OutputHandler(args, params)
	}
	poolSize := len(pool)
	Evaluation(startTime, poolSize)
}

func DirectEnum(args *Args) error {
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
	excludes := strings.Split(args.ExcHttpCodes, ",")
	scanner := bufio.NewScanner(stream)
	fmt.Println()
	for scanner.Scan() {
		entry := scanner.Text()
		url := fmt.Sprintf("http://%s.%s", entry, args.Host)
		statusCode := HttpStatusCode(url)
		code := fmt.Sprintf("%d", statusCode)
		if len(excludes) != 0 && IsInExclude(code, excludes) {
			continue
		}
		fmt.Printf(" ===[ %s.%s: %d\n", entry, args.Host, statusCode)
		counter++
	}
	if err := scanner.Err(); err != nil {
		return errors.New("scanner returns an error while reading wordlist")
	}
	Evaluation(startTime, counter)
	return nil
}
