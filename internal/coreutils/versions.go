package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PlagueByteSec/sdakit-project/v2/internal/logging"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/requests"
	"github.com/PlagueByteSec/sdakit-project/v2/internal/shared"

	"github.com/hashicorp/go-version"
)

func GetCurrentRepoVersion(client *http.Client) string {
	/*
		Request the version.txt file from GitHub and return
		the value as a string.
	*/
	_, _, responseBody, _ := requests.RequestHandlerCore(&requests.HttpRequestBase{
		HttpClient:       client,
		CustomUrl:        shared.VersionUrl,
		HttpMethod:       "GET",
		ResponseNeedBody: true,
	})
	return string(responseBody)
}

func GetCurrentLocalVersion() string {
	/*
		Read the version of the current local project instance. If an error
		occurs while trying to read version.txt, return n/a.
	*/
	cwd, err := os.Getwd()
	if err != nil {
		logging.GLogger.Log(err.Error())
		return shared.NotAvailable
	}
	versionFilePath := filepath.Join(cwd, shared.VersionFile)
	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		logging.GLogger.Log(err.Error())
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
	parseRepoVersion, err := version.NewVersion(versionRepo)
	if err != nil {
		return
	}
	parseLocalVersion, err := version.NewVersion(versionLocal)
	if err != nil {
		return
	}
	if versionRepo != versionLocal && parseLocalVersion.LessThan(parseRepoVersion) {
		fmt.Fprintf(shared.GStdout, "[*] An update is available! %s->%s\n", versionLocal, versionRepo)
	}
}
