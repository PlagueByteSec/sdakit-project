package main

import (
	"Sentinel/lib"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var (
		pool         lib.Pool
		httpClient   *http.Client
		err          error
		localVersion string
		repoVersion  string
		sigChan      chan os.Signal
	)
	args, err := lib.CliParser()
	if err != nil {
		fmt.Fprintln(lib.GStdout, err)
		lib.GStdout.Flush()
		goto exitErr
	}
	if args.Verbose {
		lib.GVerbose = true
	}
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			fmt.Fprintln(lib.GStdout, "\n\nG0oDBy3!")
			lib.SentinelExit(&pool)
		}
	}()
	httpClient, err = lib.HttpClientInit(&args)
	if err != nil {
		goto exitErr
	}
	localVersion = lib.GetCurrentLocalVersion()
	repoVersion = lib.GetCurrentRepoVersion(httpClient)
	lib.VersionCompare(repoVersion, localVersion)
	fmt.Fprintf(lib.GStdout, " ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	lib.DisplayCount = 0
	if args.NewOutputPath == "defaultPath" {
		args.NewOutputPath = lib.OutputDir
	} else {
		lib.VerbosePrint("[*] New output directory path set: %s\n", args.NewOutputPath)
	}
	if err := lib.CreateOutputDir(args.NewOutputPath); err != nil {
		goto exitErr
	}
	fmt.Fprint(lib.GStdout, "[*] Method: ")
	if len(args.WordlistPath) == 0 {
		fmt.Fprintln(lib.GStdout, "PASSIVE")
		lib.PassiveEnum(&args, httpClient)
	} else {
		fmt.Fprintln(lib.GStdout, "DIRECT")
		if err := lib.DirectEnum(&args, httpClient); err != nil {
			goto exitErr
		}
	}
	lib.SentinelExit(&pool)
exitErr:
	lib.Logger.Println(err)
	lib.Logger.Fatalf("Program execution failed")
}
