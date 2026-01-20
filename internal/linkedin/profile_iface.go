package linkedin

import (
	"context"
)

// ProfileInterface defines the interface for fetching LinkedIn profile data
type ProfileInterface interface {
	// GetProfile retrieves the user's LinkedIn profile using an access token
	GetProfile(ctx context.Context, accessToken string) (*Profile, error)

	// GetProfileByUsername retrieves a specific LinkedIn profile by username using an access token
	GetProfileByUsername(ctx context.Context, accessToken, username string) (*Profile, error)
}
