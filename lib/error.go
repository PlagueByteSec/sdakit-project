package lib

type VersionHandler struct{}

type TestVersion interface {
	HandleVersion(err error)
}

func (handler *VersionHandler) HandleVersion(err error) string {
	var version string
	if err != nil {
		version = na
	}
	return version
}

func TestVersionFail(handler VersionHandler, version *string, err error) {
	check := handler.HandleVersion(err)
	if check == na {
		*version = check
	}
}
