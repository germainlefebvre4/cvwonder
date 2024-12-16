package model

type CV struct {
	LiveReload any `yaml:"liveReload"`
	Site       struct {
		URL string `yaml:"url"`
	} `yaml:"site"`
	Company struct {
		Name string `yaml:"name"`
		Logo string `yaml:"logo"`
	} `yaml:"company"`
	Person struct {
		Name        string `yaml:"name"`
		Depiction   string `yaml:"depiction"`
		Profession  string `yaml:"profession"`
		Location    string `yaml:"location"`
		Citizenship string `yaml:"citizenship"`
		Email       string `yaml:"email"`
		Site        string `yaml:"site"`
	} `yaml:"person"`
	SocialNetworks struct {
		Github        string `yaml:"github"`
		Stackoverflow string `yaml:"stackoverflow"`
		Linkedin      string `yaml:"linkedin"`
		Twitter       string `yaml:"twitter"`
	} `yaml:"social_networks"`
	Abstract []struct {
		Tr string `yaml:"tr"`
	} `yaml:"abstract"`
	Career []struct {
		CompanyName string `yaml:"companyName"`
		Duration    string `yaml:"duration,omitempty"`
		Missions    []struct {
			Position     string   `yaml:"position"`
			Company      string   `yaml:"company"`
			Location     string   `yaml:"location"`
			Dates        string   `yaml:"dates"`
			Summary      string   `yaml:"summary"`
			Technologies []string `yaml:"technologies"`
			Description  []string `yaml:"description"`
			Project      string   `yaml:"project,omitempty"`
		} `yaml:"missions"`
	} `yaml:"career"`
	TechnicalSkills struct {
		List    []interface{} `yaml:"list"`
		Domains []struct {
			Name         string `yaml:"name"`
			Competencies []struct {
				Name  string `yaml:"name"`
				Level int    `yaml:"level"`
			} `yaml:"competencies"`
		} `yaml:"domains"`
	} `yaml:"technicalSkills"`
	SideProjects []struct {
		Name        string `yaml:"name"`
		Position    string `yaml:"position"`
		Description string `yaml:"description"`
		Link        string `yaml:"link"`
		Type        string `yaml:"type"`
		Langs       string `yaml:"langs"`
		Color       string `yaml:"color"`
	} `yaml:"sideProjects"`
	Certifications []struct {
		CompanyName       string `yaml:"companyName"`
		CertificationName string `yaml:"certificationName"`
		Issuer            string `yaml:"issuer"`
		Date              string `yaml:"date"`
		Link              string `yaml:"link"`
		Badge             string `yaml:"badge"`
	} `yaml:"certifications"`
	Education []struct {
		Name     string `yaml:"name"`
		Degree   string `yaml:"degree"`
		Location string `yaml:"location"`
		Dates    string `yaml:"dates"`
		Link     string `yaml:"link"`
	} `yaml:"education"`
}
