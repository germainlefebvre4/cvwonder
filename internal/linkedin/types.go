package linkedin

// Profile represents a LinkedIn user profile
type Profile struct {
	// Basic information
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Headline  string `json:"headline"`

	// Contact information
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`

	// Profile links
	ProfileURL string `json:"profileURL"`
	PhotoURL   string `json:"photoURL"`

	// Professional information
	Summary        string          `json:"summary"`
	Positions      []Position      `json:"positions"`
	Education      []Education     `json:"education"`
	Skills         []Skill         `json:"skills"`
	Languages      []Language      `json:"languages"`
	Projects       []Project       `json:"projects"`
	Certifications []Certification `json:"certifications"`
}

// Position represents a work position
type Position struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	CompanyLogo string `json:"companyLogo"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	IsCurrent   bool   `json:"isCurrent"`
	Description string `json:"description"`
}

// Education represents educational background
type Education struct {
	ID           string `json:"id"`
	School       string `json:"school"`
	SchoolLogo   string `json:"schoolLogo"`
	Degree       string `json:"degree"`
	FieldOfStudy string `json:"fieldOfStudy"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Location     string `json:"location"`
}

// Skill represents a professional skill
type Skill struct {
	Name string `json:"name"`
}

// Language represents a spoken language
type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

// Project represents a side project
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

// Certification represents a professional certification
type Certification struct {
	Name          string `json:"name"`
	Authority     string `json:"authority"`
	LicenseNumber string `json:"licenseNumber"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	URL           string `json:"url"`
}
