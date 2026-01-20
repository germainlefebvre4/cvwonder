package linkedin

import (
	"context"
)

// AuthInterface defines the interface for LinkedIn authentication
type AuthInterface interface {
	// Authenticate performs OAuth2 authentication with LinkedIn
	// Returns an access token on success
	Authenticate(ctx context.Context, clientID, clientSecret string) (string, error)

	// GetAuthorizationURL returns the URL for user authorization
	GetAuthorizationURL(clientID, redirectURI, state string) string

	// ExchangeCodeForToken exchanges an authorization code for an access token
	ExchangeCodeForToken(ctx context.Context, clientID, clientSecret, code, redirectURI string) (string, error)
}
