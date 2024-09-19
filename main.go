package main

import (
	"github.com/PlagueByteSec/sdakit-project/v2/cmd"
	utils "github.com/PlagueByteSec/sdakit-project/v2/internal/coreutils"
)

func main() {
	args, err := cmd.CliParser()
	if err != nil {
		utils.ProgramExit(utils.ExitParams{
			ExitCode:    -1,
			ExitMessage: "CliParser failed",
			ExitError:   err,
		})
	}
	cmd.Run(args)
}
