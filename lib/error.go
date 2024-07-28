package lib

import (
	"fmt"
	"os"
)

func GetPanic(base string, args ...interface{}) {
	message := fmt.Sprintf(base, args...)
	fmt.Println("[-] ERROR: " + message)
	os.Exit(-1)
}

func TestVersionFail(err error) string {
	var value string
	if err != nil {
		value = "n/a"
	}
	return value
}
