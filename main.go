package main

import (
	"Sentinel/lib"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DefaultOutputName(hostname string) string {
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02_15-04-05")
	outputFile := fmt.Sprintf("%s-%s.txt", formatTime, hostname)
	return outputFile
}

func CreateOutputDir() {
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			lib.GetPanic("Unable to create output directory: %s\n", outputDir)
		}
	}
}

// Pool init and preparation
func EntryPoint(args *lib.Args) {
	startTime := time.Now()
	CreateOutputDir()
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
			defaultOutput := DefaultOutputName(args.Host)
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
	args := lib.CliParser()
	EntryPoint(&args)
}
