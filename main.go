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
	for _, entry := range lib.Db {
		url := strings.Replace(entry, "HOST", args.Host, 1)
		lib.Request(pool, args.Host, url)
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

func DirectEnum(args *lib.Args) {
	var counter int
	startTime := time.Now()
	if _, err := os.Stat(args.WordlistPath); errors.Is(err, os.ErrNotExist) {
		lib.GetPanic("Could not find wordlist \"%s\" :(", args.WordlistPath)
	}
	stream, err := os.Open(args.WordlistPath)
	if err != nil {
		lib.GetPanic("%s", err)
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
		lib.GetPanic("%s", err)
	}
	Evaluation(startTime, counter)
}

func main() {
	args := lib.CliParser()
	lib.VersionCompare()
	lib.CreateOutputDir()
	fmt.Println("[*] Sending GET request to endpoints..")
	if len(args.WordlistPath) == 0 {
		PassiveEnum(&args)
	} else {
		DirectEnum(&args)
	}
}
