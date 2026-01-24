package linkedin

import (
	"fmt"
	"strings"

	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
)

// ConverterService implements the ConverterInterface
type ConverterService struct{}

// NewConverterService creates a new ConverterService
func NewConverterService() (*ConverterService, error) {
	return &ConverterService{}, nil
}

// Convert transforms a LinkedIn profile into CVWonder CV model
func (c *ConverterService) Convert(profile *Profile) (*model.CV, error) {
	logrus.Debug("Converting LinkedIn profile to CV model")

	cv := &model.CV{
		Person: model.Person{
			Name:       fmt.Sprintf("%s %s", profile.FirstName, profile.LastName),
			Profession: profile.Headline,
			Location:   profile.Location,
			Email:      profile.Email,
			Depiction:  profile.PhotoURL,
			Phone:      profile.PhoneNumber,
		},
		SocialNetworks: model.SocialNetworks{
			Linkedin: profile.ProfileURL,
		},
	}

	// Convert abstract/summary
	if profile.Summary != "" {
		cv.Abstract = strings.Split(profile.Summary, "\n")
	}

	// Convert positions to career
	cv.Career = c.convertPositions(profile.Positions)

	// Convert education
	cv.Education = c.convertEducation(profile.Education)

	// Convert skills
	cv.TechnicalSkills = c.convertSkills(profile.Skills)

	// Convert projects
	cv.SideProjects = c.convertProjects(profile.Projects)

	// Convert certifications
	cv.Certifications = c.convertCertifications(profile.Certifications)

	// Convert languages
	cv.Languages = c.convertLanguages(profile.Languages)

	logrus.Debug("Successfully converted LinkedIn profile to CV model")
	return cv, nil
}

// ConvertToYAML transforms a LinkedIn profile into YAML format
func (c *ConverterService) ConvertToYAML(profile *Profile) ([]byte, error) {
	logrus.Debug("Converting LinkedIn profile to YAML")

	cv, err := c.Convert(profile)
	if err != nil {
		return nil, fmt.Errorf("failed to convert profile: %w", err)
	}

	yamlData, err := yaml.Marshal(cv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CV to YAML: %w", err)
	}

	logrus.Debug("Successfully converted LinkedIn profile to YAML")
	return yamlData, nil
}

// convertPositions converts LinkedIn positions to career missions
func (c *ConverterService) convertPositions(positions []Position) []model.Career {
	if len(positions) == 0 {
		return []model.Career{}
	}

	// Group positions by company
	companyMap := make(map[string][]Position)
	for _, pos := range positions {
		companyMap[pos.Company] = append(companyMap[pos.Company], pos)
	}

	var career []model.Career
	for companyName, companyPositions := range companyMap {
		var missions []model.Mission

		for _, pos := range companyPositions {
			var dates string
			if pos.EndDate != "" {
				dates = fmt.Sprintf("%s - %s", pos.StartDate, pos.EndDate)
			} else {
				dates = fmt.Sprintf("%s - Present", pos.StartDate)
			}

			mission := model.Mission{
				Position: pos.Title,
				Company:  pos.Company,
				Location: pos.Location,
				Dates:    dates,
				Summary:  pos.Description,
			}

			if pos.Description != "" {
				mission.Description = strings.Split(pos.Description, "\n")
			}

			missions = append(missions, mission)
		}

		careerEntry := model.Career{
			CompanyName: companyName,
			CompanyLogo: companyPositions[0].CompanyLogo,
			Missions:    missions,
		}

		career = append(career, careerEntry)
	}

	return career
}

// convertEducation converts LinkedIn education to CV education
func (c *ConverterService) convertEducation(education []Education) []model.Education {
	var result []model.Education

	for _, edu := range education {
		dates := ""
		if edu.StartDate != "" && edu.EndDate != "" {
			dates = fmt.Sprintf("%s - %s", edu.StartDate, edu.EndDate)
		} else if edu.StartDate != "" {
			dates = edu.StartDate
		}

		degree := edu.Degree
		if edu.FieldOfStudy != "" {
			degree = fmt.Sprintf("%s in %s", edu.Degree, edu.FieldOfStudy)
		}

		result = append(result, model.Education{
			SchoolName: edu.School,
			SchoolLogo: edu.SchoolLogo,
			Degree:     degree,
			Location:   edu.Location,
			Dates:      dates,
		})
	}

	return result
}

// convertSkills converts LinkedIn skills to technical skills
func (c *ConverterService) convertSkills(skills []Skill) model.TechnicalSkills {
	if len(skills) == 0 {
		return model.TechnicalSkills{}
	}

	// Create a single domain for all skills
	var competencies []model.Competency
	for _, skill := range skills {
		competencies = append(competencies, model.Competency{
			Name:  skill.Name,
			Level: 3, // Default level
		})
	}

	return model.TechnicalSkills{
		Domains: []model.Domain{
			{
				Name:         "Technical Skills",
				Competencies: competencies,
			},
		},
	}
}

// convertProjects converts LinkedIn projects to side projects
func (c *ConverterService) convertProjects(projects []Project) []model.SideProject {
	var result []model.SideProject

	for _, proj := range projects {
		result = append(result, model.SideProject{
			Name:        proj.Name,
			Description: proj.Description,
			Link:        proj.URL,
			Type:        "project",
		})
	}

	return result
}

// convertCertifications converts LinkedIn certifications to CV certifications
func (c *ConverterService) convertCertifications(certifications []Certification) []model.Certification {
	var result []model.Certification

	for _, cert := range certifications {
		result = append(result, model.Certification{
			CertificationName: cert.Name,
			Issuer:            cert.Authority,
			Date:              cert.StartDate,
			Link:              cert.URL,
		})
	}

	return result
}

// convertLanguages converts LinkedIn languages to CV languages
func (c *ConverterService) convertLanguages(languages []Language) []model.Language {
	var result []model.Language

	for _, lang := range languages {
		result = append(result, model.Language{
			Name:  lang.Name,
			Level: lang.Proficiency,
		})
	}

	return result
}
