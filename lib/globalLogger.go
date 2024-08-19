package lib

import (
	"Sentinel/lib/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	/*
		Create the log directory if it does not exist, and use the log file name with
		the pattern <date>-sentinel.log to log all messages.
	*/
	if err := utils.CreateOutputDir(utils.LoggerOutputDir); err != nil {
		fmt.Println("[-] Failed to create output directory for global logger. No logs will be available!")
		return
	}
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	logFileName := fmt.Sprintf("%s-%s", formatTime, utils.LogFileName)
	logFilePath := filepath.Join(utils.LoggerOutputDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	utils.Glogger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
