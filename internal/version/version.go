package version

import (
	"github.com/sirupsen/logrus"
)

var CVWONDER_VERSION = "0.3.1"

func (t *VersionService) GetVersion() {
	logrus.Debug("GetVersion")
	logrus.Info(Version())
}

func Version() string {
	version := "Version: " + CVWONDER_VERSION
	return version
}
