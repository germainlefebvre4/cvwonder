package version

import (
	"log/slog"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/logger"
	"github.com/stretchr/testify/assert"
)

func NewVersionServicesTest() VersionService {
	return VersionService{}
}

func TestGetVersion(t *testing.T) {
	out, _ := logger.Capture(t, slog.LevelDebug)

	service := &VersionService{}

	assert.Equal(t, CVWONDER_VERSION, service.GetVersion())
	assert.Contains(t, out.String(), "GetVersion", "Expected debug log to contain 'GetVersion'")
}
