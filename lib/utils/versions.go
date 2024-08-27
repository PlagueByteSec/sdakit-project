package utils

import (
	"Sentinel/lib/shared"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
)

func GetCurrentLocalVersion() string {
	/*
		Read the version of the current local project instance. If an error
		occurs while trying to read version.txt, return n/a.
	*/
	cwd, err := os.Getwd()
	if err != nil {
		shared.Glogger.Println(err)
		return shared.NotAvailable
	}
	versionFilePath := filepath.Join(cwd, shared.VersionFile)
	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		shared.Glogger.Println(err)
		return shared.NotAvailable
	}
	return string(content)
}

func VersionCompare(versionRepo string, versionLocal string) {
	/*
		Compare the version of the local project instance with the version
		from the GitHub repository. If the local version is lower than the repository
		version, the user is notified that updates are available.
	*/
	if versionRepo == shared.NotAvailable ||
		versionLocal == shared.NotAvailable || versionLocal == "" {
		return
	}
	parseRepoVersion, _ := version.NewVersion(versionRepo)
	parseLocalVersion, _ := version.NewVersion(versionLocal)
	if versionRepo != versionLocal && parseLocalVersion.LessThan(parseRepoVersion) {
		fmt.Fprintf(shared.GStdout, "[*] An update is available! %s->%s\n", versionLocal, versionRepo)
	}
}
