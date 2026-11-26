package version

type VersionInterface interface {
	GetVersion() string
}

type VersionService struct{}

func NewVersionService() (VersionInterface, error) {
	return &VersionService{}, nil
}
