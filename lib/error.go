package lib

type VersionHandler struct{}

type TestVersion interface {
	HandleVersion(err error)
}

func (handler *VersionHandler) handleVersion(err error) string {
	var version string
	if err != nil {
		version = Na
	}
	return version
}

func TestVersionFail(handler VersionHandler, version *string, err error) {
	check := handler.handleVersion(err)
	if check == Na {
		*version = check
	}
}
