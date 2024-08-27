package lib

import (
	"Sentinel/lib/externs"
	"Sentinel/lib/shared"
	"Sentinel/lib/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	// make every pool at startup and open a stream writer to stdout.
	shared.PoolInit(&shared.GPoolBase)
	shared.GStdout = bufio.NewWriter(os.Stdout)
	/*
		Create the log directory if it does not exist, and use the log file name with
		the pattern <date>-sentinel.log to log all messages.
	*/
	if err := utils.CreateOutputDir(shared.LoggerOutputDir); err != nil {
		fmt.Println("[-] Failed to create output directory for global logger. No logs will be available!")
		return
	}
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	logFileName := fmt.Sprintf("%s-%s", formatTime, shared.LogFileName)
	logFilePath := filepath.Join(shared.LoggerOutputDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	shared.Glogger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

func MethodManagerInit() map[string]shared.EnumerationMethod {
	return map[string]shared.EnumerationMethod{
		shared.Passive: {
			MethodKey: shared.Passive,
			Action:    PassiveEnum,
		},
		shared.Active: {
			MethodKey: shared.Active,
			Action:    ActiveEnum,
		},
		shared.Dns: {
			MethodKey: shared.Dns,
			Action:    DnsEnum,
		},
	}
}

func ExternsManagerInit() map[string]shared.ExternsMethod {
	return map[string]shared.ExternsMethod{
		shared.RDns: {
			MethodKey: shared.RDns,
			Action:    externs.RDnsFromFile,
		},
		shared.Ping: {
			MethodKey: shared.Ping,
			Action:    externs.PingFromFile,
		},
	}
}
