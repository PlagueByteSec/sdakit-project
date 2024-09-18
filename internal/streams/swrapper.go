package streams

import (
	"fmt"
	"net/http"

	utils "github.com/PlagueByteSec/sentinel-project/v2/internal/coreutils"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/logging"
	"github.com/PlagueByteSec/sentinel-project/v2/internal/shared"
	"github.com/PlagueByteSec/sentinel-project/v2/pkg"
)

func OutputHandlerWrapper(subdomain string, client *http.Client, args *shared.Args,
	paramsSetupFiles *shared.ParamsSetupFilesBase, url string) {
	if !utils.IsPassiveEnumeration(args) {
		dotChan := make(chan struct{})
		go pkg.PrintDots(subdomain, dotChan)
		fmt.Fprintf(shared.GStdout, "\rFOUND: %s, analyzing", subdomain)
		shared.GStdout.Flush()
		OutputHandler(&shared.GStreams, client, args, *paramsSetupFiles.FileParams, url)
		close(dotChan)
	} else {
		OutputHandler(&shared.GStreams, client, args, *paramsSetupFiles.FileParams, url)
	}
}

func OpenOutputFileStreamsWrapper(filePaths *shared.FilePaths) {
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if err := OpenOutputFileStreams(&shared.GStreams, filePaths); err != nil {
		logging.GLogger.Log(err.Error())
	}
}

func OutputWrapper(ipAddrs []string, params shared.Params, streams *shared.FileStreams) {
	for _, ip := range ipAddrs {
		IpManage(params, ip, streams)
	}
	if !shared.GDisableAllOutput {
		err := WriteOutputFileStream(streams.SubdomainStream, params.FileContentSubdoms)
		if err != nil {
			streams.SubdomainStream.Close()
		}
	}
}
