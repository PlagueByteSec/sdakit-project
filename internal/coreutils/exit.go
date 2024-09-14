package utils

import (
	"fmt"
	"os"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

func SentinelExit(exitParams shared.SentinelExitParams) {
	/*
		Read the exit settings specified in SentinelExitParams and
		adjust the behavior based on those settings.
	*/
	fmt.Fprintln(shared.GStdout, exitParams.ExitMessage)
	if exitParams.ExitError != nil {
		logging.GLogger.Log(exitParams.ExitError.Error())
		fmt.Fprintln(shared.GStdout, exitParams.ExitError.Error())
	}
	shared.GStdout.Flush()
	os.Exit(exitParams.ExitCode)
}

func SentinelPanic(err error) {
	fmt.Fprintf(shared.GStdout, "\r%-50s\n", err)
	shared.GStdout.Flush()
	logging.GLogger.Log(err.Error())
	os.Exit(-1)
}
