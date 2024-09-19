package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"
)

var GLogger *Logger

type Logger struct {
	logFileStream *os.File
	logMutex      sync.Mutex
	logChannel    chan string
	logFinish     chan struct{}
}

func NewLogger() (*Logger, error) {
	currentTime := time.Now()
	formatTime := currentTime.Format("2006-01-02")
	logFileName := fmt.Sprintf("%s-%s", formatTime, shared.LogFileName)
	logFilePath := filepath.Join(shared.LoggerOutputDir, logFileName)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file stream: %s", err)
	}
	if shared.GVerbose {
		fmt.Printf("[*] log file created: %s\n", logFilePath)
	}
	return &Logger{
		logChannel:    make(chan string),
		logFinish:     make(chan struct{}),
		logFileStream: logFile,
	}, nil
}

func (logger *Logger) Log(message string) {
	logger.logChannel <- message
}

func (logger *Logger) Start() {
	go func() {
		for {
			select {
			case <-logger.logFinish:
				return
			case logMessage := <-logger.logChannel:
				logger.logMutex.Lock()
				log.SetOutput(logger.logFileStream)
				log.Println(logMessage)
				logger.logMutex.Unlock()
			}
		}
	}()
}

func (logger *Logger) Stop() {
	close(logger.logFinish)
	logger.logFileStream.Close()
}
