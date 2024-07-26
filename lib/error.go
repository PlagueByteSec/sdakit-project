package lib

import (
	"fmt"
	"os"
)

func GetPanic(base string, args ...interface{}) {
	message := fmt.Sprintf(base, args...)
	fmt.Println("ERROR: " + message)
	os.Exit(-1)
}
