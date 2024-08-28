package utils

import (
	"fmt"
	"os"

	"github.com/fhAnso/Sentinel/v1/internal/shared"
)

func SentinelExit(exitParams shared.SentinelExitParams) {
	/*
		Read the exit settings specified in SentinelExitParams and
		adjust the behavior based on those settings.
	*/
	fmt.Fprintln(shared.GStdout, exitParams.ExitMessage)
	if exitParams.ExitError != nil {
		errorMessage := fmt.Sprintf("Sentinel exit with an error: %s", exitParams.ExitError.Error())
		shared.Glogger.Println(errorMessage)
		fmt.Fprintln(shared.GStdout, errorMessage)
	}
	shared.GStdout.Flush()
	os.Exit(exitParams.ExitCode)
}

func SentinelPanic(err error) {
	fmt.Fprintf(shared.GStdout, "\r%-50s\n", err)
	shared.GStdout.Flush()
	shared.Glogger.Println(err)
	shared.Glogger.Fatalf("Program execution failed")
}
