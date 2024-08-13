package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

const logFileName = "sentinel.log"

func init() {
	if err := CreateOutputDir(LoggerOutputDir); err != nil {
		Logger.Println(err)
		return
	}
	logFilePath := filepath.Join(LoggerOutputDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
