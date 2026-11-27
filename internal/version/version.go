package version

import "log/slog"

var CVWONDER_VERSION = "dev"

// GetVersion returns the CV Wonder version string.
func (t *VersionService) GetVersion() string {
	slog.Debug("GetVersion")
	return CVWONDER_VERSION
}
