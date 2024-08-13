package main

import (
	"Sentinel/lib"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var (
		httpClient   *http.Client
		err          error
		localVersion string
		repoVersion  string
	)
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			fmt.Println("\n\nG0oDBy3!")
			os.Exit(0)
		}
	}()
	args, err := lib.CliParser()
	if err != nil {
		lib.Logger.Println(err)
		goto exitErr
	}
	httpClient, err = lib.HttpClientInit(&args)
	if err != nil {
		lib.Logger.Println(err)
		goto exitErr
	}
	localVersion = lib.GetCurrentLocalVersion()
	repoVersion = lib.GetCurrentRepoVersion(httpClient)
	fmt.Printf(" ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	lib.VersionCompare(repoVersion, localVersion)
	if err := lib.CreateOutputDir(lib.OutputDir); err != nil {
		lib.Logger.Println(err)
		goto exitErr
	}
	lib.DisplayCount = 0
	fmt.Print("[*] Method: ")
	if len(args.WordlistPath) == 0 {
		fmt.Println("PASSIVE")
		lib.PassiveEnum(&args, httpClient)
	} else {
		fmt.Println("DIRECT")
		if err := lib.DirectEnum(&args, httpClient); err != nil {
			lib.Logger.Println(err)
			goto exitErr
		}
	}
	os.Exit(0)
exitErr:
	lib.Logger.Fatalf("Program execution failed")
}
