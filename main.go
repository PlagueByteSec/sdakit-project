package main

import (
	"Sentinel/lib"
	"fmt"
	"os"
)

func main() {
	failHandler := &lib.VersionHandler{}
	localVersion := lib.GetCurrentLocalVersion(failHandler)
	args := lib.CliParser()
	fmt.Printf(" ===[ Sentinel, v%s ]===\n\n", localVersion)
	lib.VersionCompare()
	if err := lib.CreateOutputDir(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	lib.DisplayCount = 0
	httpClient := lib.ClientInit()
	if len(args.WordlistPath) == 0 {
		fmt.Println("[*] Using passive enum method")
		lib.PassiveEnum(&args, httpClient)
	} else {
		fmt.Println("[*] Using direct enum method")
		if err := lib.DirectEnum(&args, httpClient); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
