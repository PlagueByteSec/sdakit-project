package main

import (
	"github.com/fhAnso/Sentinel/v1/cmd"
	utils "github.com/fhAnso/Sentinel/v1/internal/coreutils"
)

func main() {
	args, err := cmd.CliParser()
	if err != nil {
		utils.SentinelPanic(err)
	}
	cmd.Run(args)
}
