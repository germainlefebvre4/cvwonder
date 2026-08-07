package model

// CustomField is a single author-defined {label, value} entry, used both as
// an entity-level extension point (Extensible.Custom) and inside CustomSection.
type CustomField struct {
	Label string      `yaml:"label"`
	Value interface{} `yaml:"value"`
}

// CustomSection is a whole new CV-root section the model has no dedicated
// field for, identified by a Title and holding its own list of fields.
type CustomSection struct {
	Title  string        `yaml:"title"`
	Fields []CustomField `yaml:"fields"`
}

// Extensible is embedded into model structs to give them a flat `custom:`
// YAML key (via yaml:",inline") without repeating the field declaration.
type Extensible struct {
	Custom []CustomField `yaml:"custom,omitempty"`
}

type CV struct {
	Extensible      `yaml:",inline"`
	Company         Company         `yaml:"company"`
	Person          Person          `yaml:"person"`
	SocialNetworks  SocialNetworks  `yaml:"socialNetworks"`
	Abstract        []string        `yaml:"abstract"`
	Career          []Career        `yaml:"career"`
	TechnicalSkills TechnicalSkills `yaml:"technicalSkills"`
	SideProjects    []SideProject   `yaml:"sideProjects"`
	Certifications  []Certification `yaml:"certifications"`
	Languages       []Language      `yaml:"languages"`
	Education       []Education     `yaml:"education"`
	References      []Reference     `yaml:"references"`
	CustomSections  []CustomSection `yaml:"customSections,omitempty"`
}

type Reference struct {
	Extensible     `yaml:",inline"`
	Name           string         `yaml:"name"`
	Position       string         `yaml:"position"`
	Company        string         `yaml:"company"`
	Date           string         `yaml:"date"`
	Url            string         `yaml:"url"`
	SocialNetworks SocialNetworks `yaml:"socialNetworks"`
	Description    string         `yaml:"description"`
}

type Company struct {
	Extensible `yaml:",inline"`
	Name       string `yaml:"name"`
	Logo       string `yaml:"logo"`
}

type Person struct {
	Extensible  `yaml:",inline"`
	Name        string     `yaml:"name"`
	Depiction   string     `yaml:"depiction"`
	Profession  string     `yaml:"profession"`
	Location    string     `yaml:"location"`
	Citizenship string     `yaml:"citizenship"`
	Email       string     `yaml:"email"`
	Site        string     `yaml:"site"`
	Phone       string     `yaml:"phone"`
	Experience  Experience `yaml:"experience,omitempty"`
}

type Experience struct {
	Extensible `yaml:",inline"`
	Years      int `yaml:"years,omitempty"`
	Since      int `yaml:"since,omitempty"`
}

type SocialNetworks struct {
	Extensible    `yaml:",inline"`
	Github        string `yaml:"github,omitempty"`
	Stackoverflow string `yaml:"stackoverflow,omitempty"`
	Linkedin      string `yaml:"linkedin,omitempty"`
	Twitter       string `yaml:"twitter,omitempty"`
	Bluesky       string `yaml:"bluesky,omitempty"`
}

type Career struct {
	Extensible  `yaml:",inline"`
	CompanyName string    `yaml:"companyName"`
	CompanyLogo string    `yaml:"companyLogo"`
	Duration    string    `yaml:"duration,omitempty"`
	Missions    []Mission `yaml:"missions"`
}

type Mission struct {
	Extensible   `yaml:",inline"`
	Position     string   `yaml:"position"`
	Company      string   `yaml:"company"`
	CompanyLogo  string   `yaml:"companyLogo,omitempty"`
	Location     string   `yaml:"location"`
	Dates        string   `yaml:"dates"`
	Summary      string   `yaml:"summary"`
	Technologies []string `yaml:"technologies"`
	Description  []string `yaml:"description"`
	Project      string   `yaml:"project,omitempty"`
}

type TechnicalSkills struct {
	Extensible `yaml:",inline"`
	Domains    []Domain `yaml:"domains"`
}

type Domain struct {
	Extensible   `yaml:",inline"`
	Name         string       `yaml:"name"`
	Competencies []Competency `yaml:"competencies"`
}

type Competency struct {
	Extensible `yaml:",inline"`
	Name       string `yaml:"name"`
	Level      int    `yaml:"level"`
}

type SideProject struct {
	Extensible  `yaml:",inline"`
	Name        string `yaml:"name"`
	Position    string `yaml:"position"`
	Description string `yaml:"description"`
	Link        string `yaml:"link"`
	Type        string `yaml:"type"`
	Langs       string `yaml:"langs"`
	Color       string `yaml:"color"`
}

type Certification struct {
	Extensible        `yaml:",inline"`
	CompanyName       string `yaml:"companyName"`
	CertificationName string `yaml:"certificationName"`
	Issuer            string `yaml:"issuer"`
	Date              string `yaml:"date"`
	Link              string `yaml:"link"`
	Badge             string `yaml:"badge"`
}

type Education struct {
	Extensible `yaml:",inline"`
	SchoolName string `yaml:"schoolName"`
	SchoolLogo string `yaml:"schoolLogo"`
	Degree     string `yaml:"degree"`
	Location   string `yaml:"location"`
	Dates      string `yaml:"dates"`
	Link       string `yaml:"link"`
}

type Language struct {
	Extensible `yaml:",inline"`
	Name       string `yaml:"name"`
	Level      string `yaml:"level"`
}
