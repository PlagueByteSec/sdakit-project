package lib

import (
	"fmt"
	"os"
)

type VersionHandler struct{}

type TestVersion interface {
	HandleVersion(err error)
}

func (handler *VersionHandler) HandleVersion(err error) string {
	var version string
	if err != nil {
		version = "n/a"
	}
	return version
}

func TestVersionFail(handler VersionHandler, version *string, err error) {
	check := handler.HandleVersion(err)
	if check == "n/a" {
		*version = handler.HandleVersion(err)
	}
}

func GetPanic(base string, args ...interface{}) {
	message := fmt.Sprintf(base, args...)
	fmt.Println("[-] ERROR: " + message)
	os.Exit(-1)
}
