package themes

type ThemesInterface interface {
	List() error
	Install(theme string, force bool) error
	Create(theme string) error
	Verify(theme string) error
	Screenshot(theme string) error
}

type ThemesService struct{}

func NewThemesService() (ThemesInterface, error) {
	return &ThemesService{}, nil
}
