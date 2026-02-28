package linkedin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthService(t *testing.T) {
	service, err := NewAuthService()
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.httpClient)
}

func TestGetAuthorizationURL(t *testing.T) {
	service, _ := NewAuthService()

	tests := []struct {
		name         string
		clientID     string
		redirectURI  string
		state        string
		wantContains []string
	}{
		{
			name:        "Should generate valid authorization URL",
			clientID:    "test-client-id",
			redirectURI: "http://localhost:8080/callback",
			state:       "test-state",
			wantContains: []string{
				"https://www.linkedin.com/oauth/v2/authorization",
				"response_type=code",
				"client_id=test-client-id",
				"redirect_uri=http",
				"state=test-state",
				"scope=openid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := service.GetAuthorizationURL(tt.clientID, tt.redirectURI, tt.state)

			for _, contains := range tt.wantContains {
				assert.Contains(t, url, contains)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	service, _ := NewAuthService()
	ctx := context.Background()

	// This method should return an error as it requires interactive authentication
	token, err := service.Authenticate(ctx, "client-id", "client-secret")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "interactive authentication required")
}
