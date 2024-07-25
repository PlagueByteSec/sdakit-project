package lib

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

// Result pool
type Pool map[string]struct{}

// Add new entry to result pool
func (pool *Pool) AddEntry(entry string) {
	(*pool)[entry] = struct{}{}
}

// Verify that the entry does not exist in Pool
func (pool *Pool) ContainsEntry(entry string) bool {
	_, exists := (*pool)[entry]
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
