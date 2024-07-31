package lib

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
		*version = check
	}
}
