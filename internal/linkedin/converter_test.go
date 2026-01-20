package linkedin

import (
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewConverterService(t *testing.T) {
	service, err := NewConverterService()
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestConvert(t *testing.T) {
	service, _ := NewConverterService()

	tests := []struct {
		name     string
		profile  *Profile
		wantErr  bool
		validate func(t *testing.T, cv *model.CV)
	}{
		{
			name: "Should convert basic profile",
			profile: &Profile{
				FirstName: "John",
				LastName:  "Doe",
				Headline:  "Software Engineer",
				Email:     "john.doe@example.com",
				Location:  "San Francisco, CA",
			},
			wantErr: false,
			validate: func(t *testing.T, cv *model.CV) {
				assert.Equal(t, "John Doe", cv.Person.Name)
				assert.Equal(t, "Software Engineer", cv.Person.Profession)
				assert.Equal(t, "john.doe@example.com", cv.Person.Email)
				assert.Equal(t, "San Francisco, CA", cv.Person.Location)
			},
		},
		{
			name: "Should convert profile with positions",
			profile: &Profile{
				FirstName: "Jane",
				LastName:  "Smith",
				Positions: []Position{
					{
						Title:       "Senior Developer",
						Company:     "Tech Corp",
						Location:    "New York, NY",
						StartDate:   "2020-01",
						EndDate:     "2023-12",
						Description: "Developed applications",
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, cv *model.CV) {
				assert.Equal(t, "Jane Smith", cv.Person.Name)
				assert.Len(t, cv.Career, 1)
				assert.Equal(t, "Tech Corp", cv.Career[0].CompanyName)
				assert.Len(t, cv.Career[0].Missions, 1)
				assert.Equal(t, "Senior Developer", cv.Career[0].Missions[0].Position)
			},
		},
		{
			name: "Should convert profile with education",
			profile: &Profile{
				FirstName: "Bob",
				LastName:  "Johnson",
				Education: []Education{
					{
						School:       "MIT",
						Degree:       "Bachelor of Science",
						FieldOfStudy: "Computer Science",
						StartDate:    "2015",
						EndDate:      "2019",
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, cv *model.CV) {
				assert.Len(t, cv.Education, 1)
				assert.Equal(t, "MIT", cv.Education[0].SchoolName)
				assert.Contains(t, cv.Education[0].Degree, "Computer Science")
			},
		},
		{
			name: "Should convert profile with skills",
			profile: &Profile{
				FirstName: "Alice",
				LastName:  "Williams",
				Skills: []Skill{
					{Name: "Go"},
					{Name: "Python"},
					{Name: "Docker"},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, cv *model.CV) {
				assert.Len(t, cv.TechnicalSkills.Domains, 1)
				assert.Len(t, cv.TechnicalSkills.Domains[0].Competencies, 3)
				assert.Equal(t, "Go", cv.TechnicalSkills.Domains[0].Competencies[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv, err := service.Convert(tt.profile)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cv)
				if tt.validate != nil {
					tt.validate(t, cv)
				}
			}
		})
	}
}

func TestConvertToYAML(t *testing.T) {
	service, _ := NewConverterService()

	profile := &Profile{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}

	yamlData, err := service.ConvertToYAML(profile)

	assert.NoError(t, err)
	assert.NotEmpty(t, yamlData)
	assert.Contains(t, string(yamlData), "person:")
	assert.Contains(t, string(yamlData), "name: Test User")
	assert.Contains(t, string(yamlData), "email: test@example.com")
}

func TestConvertPositions(t *testing.T) {
	service, _ := NewConverterService()

	positions := []Position{
		{
			Title:     "Engineer",
			Company:   "Company A",
			StartDate: "2020-01",
			EndDate:   "2021-12",
		},
		{
			Title:     "Senior Engineer",
			Company:   "Company A",
			StartDate: "2022-01",
			IsCurrent: true,
		},
	}

	career := service.convertPositions(positions)

	assert.Len(t, career, 1)
	assert.Equal(t, "Company A", career[0].CompanyName)
	assert.Len(t, career[0].Missions, 2)
}

func TestConvertEducation(t *testing.T) {
	service, _ := NewConverterService()

	education := []Education{
		{
			School:       "Test University",
			Degree:       "Master",
			FieldOfStudy: "Engineering",
			StartDate:    "2015",
			EndDate:      "2017",
		},
	}

	result := service.convertEducation(education)

	assert.Len(t, result, 1)
	assert.Equal(t, "Test University", result[0].SchoolName)
	assert.Contains(t, result[0].Degree, "Engineering")
}

func TestConvertSkills(t *testing.T) {
	service, _ := NewConverterService()

	skills := []Skill{
		{Name: "Skill1"},
		{Name: "Skill2"},
	}

	result := service.convertSkills(skills)

	assert.Len(t, result.Domains, 1)
	assert.Len(t, result.Domains[0].Competencies, 2)
	assert.Equal(t, "Skill1", result.Domains[0].Competencies[0].Name)
}

func TestConvertProjects(t *testing.T) {
	service, _ := NewConverterService()

	projects := []Project{
		{
			Name:        "Project A",
			Description: "Test project",
			URL:         "https://example.com",
		},
	}

	result := service.convertProjects(projects)

	assert.Len(t, result, 1)
	assert.Equal(t, "Project A", result[0].Name)
	assert.Equal(t, "https://example.com", result[0].Link)
}

func TestConvertCertifications(t *testing.T) {
	service, _ := NewConverterService()

	certifications := []Certification{
		{
			Name:      "Test Cert",
			Authority: "Test Authority",
			StartDate: "2023",
			URL:       "https://cert.example.com",
		},
	}

	result := service.convertCertifications(certifications)

	assert.Len(t, result, 1)
	assert.Equal(t, "Test Cert", result[0].CertificationName)
	assert.Equal(t, "Test Authority", result[0].Issuer)
}

func TestConvertLanguages(t *testing.T) {
	service, _ := NewConverterService()

	languages := []Language{
		{
			Name:        "English",
			Proficiency: "Native",
		},
	}

	result := service.convertLanguages(languages)

	assert.Len(t, result, 1)
	assert.Equal(t, "English", result[0].Name)
	assert.Equal(t, "Native", result[0].Level)
}
