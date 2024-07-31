package main

import (
	"Sentinel/lib"
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
	fmt.Printf("\n[*] %d subdomains obtained. Finished in %s\n", count, duration)
}

// Pool init and preparation
func PassiveEnum(args *lib.Args) {
	startTime := time.Now()
	pool := make(lib.Pool)
	fmt.Println("[*] Formatting db entries..")
	endpoints := lib.EditDbEntries(args.Host)
	fmt.Println("[*] Sending GET request to endpoints..")
	for idx := 0; idx < len(endpoints); idx++ {
		if err := lib.Request(pool, args.Host, endpoints[idx]); err != nil {
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
		var filePath string
		if args.OutFile == "default" {
			defaultOutput := lib.DefaultOutputName(args.Host)
			filePath = filepath.Join("output", defaultOutput)
		} else {
			filePath = args.OutFile
		}
		params := lib.Params{
			FilePath:    filePath,
			FileContent: result,
			Result:      result,
			Hostname:    args.Host,
		}
		lib.OutputWriter(*args, lib.File, params)
		lib.OutputWriter(*args, lib.Stdout, params)
	}
	poolSize := len(pool)
	Evaluation(startTime, poolSize)
}

func DirectEnum(args *lib.Args) error {
	var counter int
	startTime := time.Now()
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not find wordlist: %s", args.WordlistPath)
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
		statusCode := lib.HttpStatusCode(url)
		code := fmt.Sprintf("%d", statusCode)
		if len(excludes) != 0 && lib.IsInExclude(code, excludes) {
			continue
		}
		fmt.Printf(" ===[ %s.%s: %d\n", entry, args.Host, statusCode)
		counter++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner returns an error while reading wordlist")
	}
	Evaluation(startTime, counter)
	return nil
}

func main() {
	failHandler := &lib.VersionHandler{}
	localVersion := lib.GetCurrentLocalVersion(failHandler)
	fmt.Printf(" ===[ Sentinel, v%s ]===\n\n", localVersion)
	args := lib.CliParser()
	lib.VersionCompare()
	if err := lib.CreateOutputDir(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if len(args.WordlistPath) == 0 {
		fmt.Println("[*] Using passive enum method")
		PassiveEnum(&args)
	} else {
		fmt.Println("[*] Using direct enum method")
		if err := DirectEnum(&args); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
