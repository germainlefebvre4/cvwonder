package version

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func NewVersionServicesTest() VersionService {
	return VersionService{}
}

func TestGetVersion(t *testing.T) {
	// Set up a buffer to capture log output
	var logBuffer bytes.Buffer
	logrus.SetOutput(&logBuffer)
	logrus.SetLevel(logrus.DebugLevel)

	// Create an instance of VersionService
	service := &VersionService{}

	// Call GetVersion
	service.GetVersion()

	// Check if the debug log contains "GetVersion"
	assert.Contains(t, logBuffer.String(), "GetVersion", "Expected debug log to contain 'GetVersion'")

	// Check if the info log contains the version string
	expectedVersionLog := CVWONDER_VERSION
	assert.Contains(t, logBuffer.String(), expectedVersionLog, "Expected info log to contain the version string")
}
