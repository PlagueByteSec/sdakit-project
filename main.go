package main

import (
	"Sentinel/lib"
	"fmt"
)

func main() {
	args, err := lib.CliParser()
	if err != nil {
		lib.Logger.Println(err)
		return
	}
	httpClient, err := lib.HttpClientInit(&args)
	if err != nil {
		lib.Logger.Println(err)
		return
	}
	localVersion := lib.GetCurrentLocalVersion()
	repoVersion := lib.GetCurrentRepoVersion(httpClient)
	fmt.Printf(" ===[ Sentinel, Version: %s ]===\n\n", localVersion)
	lib.VersionCompare(repoVersion, localVersion)
	if err := lib.CreateOutputDir(lib.OutputDir); err != nil {
		lib.Logger.Println(err)
		return
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
			return
		}
	}
}
