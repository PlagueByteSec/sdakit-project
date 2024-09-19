package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	utils "github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

func init() {
	// make every pool at startup and open a stream writer to stdout.
	shared.PoolsInit(&shared.GPoolBase)
	shared.GStdout = bufio.NewWriter(os.Stdout)
	/*
		Create the log directory if it does not exist, and use the log file name with
		the pattern <date>-sentinel.log to log all messages.
	*/
	if err := pkg.CreateOutputDir(shared.LoggerOutputDir); err != nil {
		fmt.Println("[-] Failed to create output directory for global logger. No logs will be available!")
		return
	}
	var err error
	logging.GLogger, err = logging.NewLogger()
	if err != nil {
		fmt.Printf("[-] Failed to initialize logger: %s\n", err)
		return
	}
	logging.GLogger.Start()
}

func MethodManagerInit() map[string]shared.EnumerationMethod {
	return map[string]shared.EnumerationMethod{
		shared.Passive: {
			MethodKey: shared.Passive,
			Action:    PassiveEnum,
		},
		shared.Active: {
			MethodKey: shared.Active,
			Action:    DirectEnum,
		},
		shared.Dns: {
			MethodKey: shared.Dns,
			Action:    DnsEnum,
		},
		shared.VHost: {
			MethodKey: shared.VHost,
			Action:    VHostEnum,
		},
	}
}

func ValidsManagerInit() map[string]shared.ExternsMethod {
	return map[string]shared.ExternsMethod{
		shared.RDns: {
			MethodKey: shared.RDns,
			Action:    RDnsFromFile,
		},
		shared.Ping: {
			MethodKey: shared.Ping,
			Action:    PingFromFile,
		},
		shared.HeaderAnalysis: {
			MethodKey: shared.HeaderAnalysis,
			Action:    AnalyseHttpHeaderSingle,
		},
	}
}

func InterruptListenerStart() {
	// wait for interrupt signal and cancel execution
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			utils.ProgramExit(utils.ExitParams{
				ExitCode:    0,
				ExitMessage: "\n\nG0oDBy3!",
				ExitError:   nil,
			})
		}
	}()
}
