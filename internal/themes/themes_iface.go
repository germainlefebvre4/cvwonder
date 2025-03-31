package themes

type ThemesInterface interface {
	List()
	Install(theme string)
	Create(theme string)
	Verify(theme string)
}

type ThemesService struct{}

func NewThemesService() (ThemesInterface, error) {
	return &ThemesService{}, nil
}
