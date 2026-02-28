package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	linkedinAuthURL  = "https://www.linkedin.com/oauth/v2/authorization"
	linkedinTokenURL = "https://www.linkedin.com/oauth/v2/accessToken"
)

// AuthService implements the AuthInterface
type AuthService struct {
	httpClient *http.Client
}

// NewAuthService creates a new AuthService
func NewAuthService() (*AuthService, error) {
	return &AuthService{
		httpClient: &http.Client{},
	}, nil
}

// GetAuthorizationURL returns the URL for user authorization
func (a *AuthService) GetAuthorizationURL(clientID, redirectURI, state string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("scope", "openid profile email w_member_social")

	return fmt.Sprintf("%s?%s", linkedinAuthURL, params.Encode())
}

// ExchangeCodeForToken exchanges an authorization code for an access token
func (a *AuthService) ExchangeCodeForToken(ctx context.Context, clientID, clientSecret, code, redirectURI string) (string, error) {
	logrus.Debug("Exchanging authorization code for access token")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", linkedinTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	logrus.Debug("Successfully obtained access token")
	return tokenResponse.AccessToken, nil
}

// Authenticate performs OAuth2 authentication with LinkedIn
// This method is for demonstration - in practice, OAuth2 requires user interaction
func (a *AuthService) Authenticate(ctx context.Context, clientID, clientSecret string) (string, error) {
	return "", fmt.Errorf("interactive authentication required: use GetAuthorizationURL and ExchangeCodeForToken")
}
