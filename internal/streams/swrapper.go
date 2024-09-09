package streams

import "github.com/PlagueByteSec/Sentinel/v1/internal/shared"

func OpenOutputFileStreamsWrapper(filePaths *shared.FilePaths) {
	/*
		Specify the name and path for each output file. If all settings are configured, open
		separate file streams for each category (Subdomains, IPv4 addresses, and IPv6 addresses).
	*/
	if err := OpenOutputFileStreams(&shared.GStreams, filePaths); err != nil {
		shared.Glogger.Println(err)
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
