package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	linkedinAPIBaseURL = "https://api.linkedin.com/v2"
)

// ProfileService implements the ProfileInterface
type ProfileService struct {
	httpClient *http.Client
}

// NewProfileService creates a new ProfileService
func NewProfileService() (*ProfileService, error) {
	return &ProfileService{
		httpClient: &http.Client{},
	}, nil
}

// GetProfile retrieves the user's LinkedIn profile using an access token
func (p *ProfileService) GetProfile(ctx context.Context, accessToken string) (*Profile, error) {
	logrus.Debug("Fetching LinkedIn profile")

	profile := &Profile{}

	// Fetch basic profile information
	if err := p.fetchBasicProfile(ctx, accessToken, profile, ""); err != nil {
		return nil, fmt.Errorf("failed to fetch basic profile: %w", err)
	}

	// Fetch email
	if err := p.fetchEmail(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch email: ", err)
	}

	// Fetch positions (work experience)
	if err := p.fetchPositions(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch positions: ", err)
	}

	// Fetch education
	if err := p.fetchEducation(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch education: ", err)
	}

	// Fetch skills
	if err := p.fetchSkills(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch skills: ", err)
	}

	logrus.Debug("Successfully fetched LinkedIn profile")
	return profile, nil
}

// GetProfileByUsername retrieves a specific LinkedIn profile by username using an access token
func (p *ProfileService) GetProfileByUsername(ctx context.Context, accessToken, username string) (*Profile, error) {
	logrus.Debugf("Fetching LinkedIn profile for username: %s", username)

	profile := &Profile{}

	// Fetch basic profile information
	if err := p.fetchBasicProfile(ctx, accessToken, profile, username); err != nil {
		return nil, fmt.Errorf("failed to fetch basic profile: %w", err)
	}

	// Fetch email
	if err := p.fetchEmail(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch email: ", err)
	}

	// Fetch positions (work experience)
	if err := p.fetchPositions(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch positions: ", err)
	}

	// Fetch education
	if err := p.fetchEducation(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch education: ", err)
	}

	// Fetch skills
	if err := p.fetchSkills(ctx, accessToken, profile); err != nil {
		logrus.Warn("Failed to fetch skills: ", err)
	}

	logrus.Debug("Successfully fetched LinkedIn profile")
	return profile, nil
}

// fetchBasicProfile fetches basic profile information
func (p *ProfileService) fetchBasicProfile(ctx context.Context, accessToken string, profile *Profile, username string) error {
	var url string
	if username != "" {
		// Fetch profile by vanity name (username)
		url = fmt.Sprintf("%s/people/(vanityName:%s)", linkedinAPIBaseURL, username)
	} else {
		// Fetch current user's profile
		url = fmt.Sprintf("%s/me", linkedinAPIBaseURL)
	}

	var response struct {
		ID        string `json:"id"`
		FirstName struct {
			Localized map[string]string `json:"localized"`
		} `json:"firstName"`
		LastName struct {
			Localized map[string]string `json:"localized"`
		} `json:"lastName"`
		Headline struct {
			Localized map[string]string `json:"localized"`
		} `json:"headline"`
		ProfilePicture struct {
			DisplayImage string `json:"displayImage~"`
		} `json:"profilePicture"`
	}

	if err := p.makeRequest(ctx, url, accessToken, &response); err != nil {
		return err
	}

	profile.ID = response.ID

	// Extract localized strings (prefer en_US)
	for _, name := range response.FirstName.Localized {
		profile.FirstName = name
		break
	}
	for _, name := range response.LastName.Localized {
		profile.LastName = name
		break
	}
	for _, headline := range response.Headline.Localized {
		profile.Headline = headline
		break
	}

	return nil
}

// fetchEmail fetches the user's email address
func (p *ProfileService) fetchEmail(ctx context.Context, accessToken string, profile *Profile) error {
	url := fmt.Sprintf("%s/emailAddress?q=members&projection=(elements*(handle~))", linkedinAPIBaseURL)

	var response struct {
		Elements []struct {
			Handle struct {
				EmailAddress string `json:"emailAddress"`
			} `json:"handle~"`
		} `json:"elements"`
	}

	if err := p.makeRequest(ctx, url, accessToken, &response); err != nil {
		return err
	}

	if len(response.Elements) > 0 {
		profile.Email = response.Elements[0].Handle.EmailAddress
	}

	return nil
}

// fetchPositions fetches work experience
func (p *ProfileService) fetchPositions(ctx context.Context, accessToken string, profile *Profile) error {
	// Note: LinkedIn API v2 requires specific permissions for positions
	// This is a simplified implementation
	url := fmt.Sprintf("%s/positions", linkedinAPIBaseURL)

	var response struct {
		Elements []struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			CompanyName string `json:"companyName"`
			Location    string `json:"location"`
			Description string `json:"description"`
			StartDate   struct {
				Year  int `json:"year"`
				Month int `json:"month"`
			} `json:"startDate"`
			EndDate struct {
				Year  int `json:"year"`
				Month int `json:"month"`
			} `json:"endDate"`
		} `json:"elements"`
	}

	if err := p.makeRequest(ctx, url, accessToken, &response); err != nil {
		return err
	}

	for _, elem := range response.Elements {
		position := Position{
			ID:          elem.ID,
			Title:       elem.Title,
			Company:     elem.CompanyName,
			Location:    elem.Location,
			Description: elem.Description,
			StartDate:   fmt.Sprintf("%04d-%02d", elem.StartDate.Year, elem.StartDate.Month),
		}

		if elem.EndDate.Year > 0 {
			position.EndDate = fmt.Sprintf("%04d-%02d", elem.EndDate.Year, elem.EndDate.Month)
			position.IsCurrent = false
		} else {
			position.IsCurrent = true
		}

		profile.Positions = append(profile.Positions, position)
	}

	return nil
}

// fetchEducation fetches educational background
func (p *ProfileService) fetchEducation(ctx context.Context, accessToken string, profile *Profile) error {
	// Note: LinkedIn API v2 requires specific permissions for education
	url := fmt.Sprintf("%s/educations", linkedinAPIBaseURL)

	var response struct {
		Elements []struct {
			ID           string `json:"id"`
			SchoolName   string `json:"schoolName"`
			DegreeName   string `json:"degreeName"`
			FieldOfStudy string `json:"fieldOfStudy"`
			StartDate    struct {
				Year int `json:"year"`
			} `json:"startDate"`
			EndDate struct {
				Year int `json:"year"`
			} `json:"endDate"`
		} `json:"elements"`
	}

	if err := p.makeRequest(ctx, url, accessToken, &response); err != nil {
		return err
	}

	for _, elem := range response.Elements {
		education := Education{
			ID:           elem.ID,
			School:       elem.SchoolName,
			Degree:       elem.DegreeName,
			FieldOfStudy: elem.FieldOfStudy,
			StartDate:    fmt.Sprintf("%04d", elem.StartDate.Year),
			EndDate:      fmt.Sprintf("%04d", elem.EndDate.Year),
		}

		profile.Education = append(profile.Education, education)
	}

	return nil
}

// fetchSkills fetches professional skills
func (p *ProfileService) fetchSkills(ctx context.Context, accessToken string, profile *Profile) error {
	// Note: LinkedIn API v2 requires specific permissions for skills
	url := fmt.Sprintf("%s/skills", linkedinAPIBaseURL)

	var response struct {
		Elements []struct {
			Name string `json:"name"`
		} `json:"elements"`
	}

	if err := p.makeRequest(ctx, url, accessToken, &response); err != nil {
		return err
	}

	for _, elem := range response.Elements {
		skill := Skill{
			Name: elem.Name,
		}
		profile.Skills = append(profile.Skills, skill)
	}

	return nil
}

// makeRequest is a helper function to make HTTP requests to LinkedIn API
func (p *ProfileService) makeRequest(ctx context.Context, url, accessToken string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
