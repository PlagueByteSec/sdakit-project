package utils

import (
	"fmt"
	"os"

	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
)

type ExitParams struct {
	ExitCode    int
	ExitMessage string
	ExitError   error
}

func cleanEmptyFiles() {
	files := []string{
		shared.GCurrentIPv4Filename,
		shared.GCurrentIPv6Filename,
	}
	for idx := 0; idx < len(files); idx++ {
		currentFile := files[idx]
		checkFile, err := os.Stat(currentFile)
		if err != nil {
			logging.GLogger.Log(err.Error())
			continue
		}
		if checkFile.Size() == 0 {
			os.Remove(currentFile)
		}
	}
}

func ProgramExit(exitParams ExitParams) {
	/*
		Read the exit settings and
		adjust the behavior based on those settings.
	*/
	fmt.Fprintln(shared.GStdout, exitParams.ExitMessage)
	if exitParams.ExitError != nil {
		logging.GLogger.Log(exitParams.ExitError.Error())
		fmt.Fprintf(shared.GStdout, "\r%-50s\n", exitParams.ExitError.Error())
	}
	cleanEmptyFiles()
	shared.GStdout.Flush()
	os.Exit(exitParams.ExitCode)
}
