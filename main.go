package main

import (
	"Sentinel/lib"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Pool init and preparation
func EntryPoint(args *lib.Args) {
	startTime := time.Now()
	lib.CreateOutputDir()
	pool := make(lib.Pool)
	fmt.Println("[*] Sending GET request to endpoints..")
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
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("\n[*] %d subdomains obtained. Finished in %s\n",
		len(pool), duration)
}

func main() {
	lib.VersionCompare()
	os.Exit(0)
	args := lib.CliParser()
	EntryPoint(&args)
}
