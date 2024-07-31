package lib

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type OutputType int

const (
	Stdout OutputType = iota
	File   OutputType = iota
)

type Params struct {
	FilePath    string
	FileContent string
	Result      string
	Hostname    string
}

func FileWriteResults(param Params) error {
	hasExt := strings.HasSuffix(param.FilePath, ".txt")
	if !hasExt {
		param.FilePath = param.FilePath + ".txt"
	}
	stream, err := os.OpenFile(param.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.New("failed to open output file stream")
	}
	defer stream.Close()
	if _, err = stream.WriteString(param.FileContent + "\n"); err != nil {
		return errors.New("output write operation failed")
	}
	return nil
}

func StdoutWriteResults(args *Args, params Params) {
	consoleOutput := fmt.Sprintf(" ===[ %s", params.Result)
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Result)
		httpStatusCode := HttpStatusCode(url)
		consoleOutput = fmt.Sprintf("%s, HTTP Status Code: %d", consoleOutput, httpStatusCode)
	}
	fmt.Println(consoleOutput)
}

func OutputWriter(args Args, outputType OutputType, params Params) {
	switch outputType {
	case Stdout:
		StdoutWriteResults(&args, params)
	case File:
		if err := FileWriteResults(params); err != nil {
			fmt.Println(err)
		}
	}
}
