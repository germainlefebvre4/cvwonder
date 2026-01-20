package linkedin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProfileService(t *testing.T) {
	service, err := NewProfileService()
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.httpClient)
}

func TestProfileServiceStructure(t *testing.T) {
	service, err := NewProfileService()

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Implements(t, (*ProfileInterface)(nil), service)
}

func TestProfileServiceImplementsGetProfileByUsername(t *testing.T) {
	service, err := NewProfileService()

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Verify that the service implements the interface including GetProfileByUsername
	var _ ProfileInterface = service
}
