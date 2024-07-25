package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Result pool
type Pool map[string]struct{}

// Add new entry to result pool
func (pool Pool) AddEntry(entry string) {
	pool[entry] = struct{}{}
}

// Verify that the entry does not exist in Pool
func (pool Pool) ContainsEntry(entry string) bool {
	_, exists := pool[entry]
	return exists
}

// Identify OS and adjust path
func DefaultPath() string {
	var filePath string
	if runtime.GOOS == "windows" {
		filePath = ".\\assets\\URL.txt"
	} else if runtime.GOOS == "linux" {
		filePath = "../assets/URL.txt"
	} else {
		fmt.Println("Your OS is currently NOT supported. Sorry :(")
		os.Exit(-1)
	}
	return filePath
}

// Make sure URL.txt exist
func FileExist(filePath string) bool {
	var exist bool
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		exist = false
	} else {
		exist = true
	}
	return exist
}

// Read database file (URL.txt) and process the entries
func Manager(host string) {
	filePath := DefaultPath()
	if FileExist(filePath) {
		pool := make(Pool)
		stream, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("ERROR: failed to open file: %s\n%s\n", filePath, err)
			os.Exit(-1)
		}
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") {
				url := strings.Replace(line, "HOST", host, 1)
				Request(pool, host, url)
			}
		}
	} else {
		fmt.Printf("ERROR: file not found in default path: %s\n", filePath)
		os.Exit(-1)
	}
}
