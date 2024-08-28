package main

import (
	"github.com/fhAnso/Sentinel/v1/cmd"
	"github.com/fhAnso/Sentinel/v1/internal/utils"
)

func main() {
	args, err := cmd.CliParser()
	if err != nil {
		utils.SentinelPanic(err)
	}
	cmd.Run(args)
}
