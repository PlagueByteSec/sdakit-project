package main

import (
	"Sentinel/lib"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Args struct {
	host        string
	outFile     string
	httpCode    bool
	pingResults bool
}

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
func EntryPoint(args Args) {
	CreateOutputDir()
	pool := make(lib.Pool)
	for _, entry := range lib.Db {
		url := strings.Replace(entry, "HOST", args.host, 1)
		lib.Request(pool, args.host, url)
	}
	for result := range pool {
		fmt.Println(args.host, ": ", result)
		if args.outFile == "default" {
			defaultOutput := DefaultOutputName(args.host)
			filePath := filepath.Join("output", defaultOutput)
			lib.WriteOutput(filePath, result)
		}
	}
}

func main() {
	host := flag.String("t", "", "Target host")
	outFile := flag.String("o", "default", "Output file")
	httpCode := flag.Bool("c", false, "Get HTTP status code of each entry")
	pingResults := flag.Bool("p", false, "Send ICMP packet to each entry")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println(lib.Help)
		os.Exit(-1)
	}
	args := Args{
		host:        *host,
		outFile:     *outFile,
		httpCode:    *httpCode,
		pingResults: *pingResults,
	}
	EntryPoint(args)
}
