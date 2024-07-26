package lib

import (
	"fmt"
	"os"
)

func GetPanic(base string, args ...interface{}) {
	message := fmt.Sprintf(base, args...)
	fmt.Println(message)
	os.Exit(-1)
}
