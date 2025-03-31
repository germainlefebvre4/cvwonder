package version

type VersionInterface interface {
	GetVersion()
}

type VersionService struct{}

func NewVersionService() (VersionInterface, error) {
	return &VersionService{}, nil
}
