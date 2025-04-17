package version

import (
	"github.com/sirupsen/logrus"
)

var CVWONDER_VERSION = "dev"

func (t *VersionService) GetVersion() {
	logrus.Debug("GetVersion")
	logrus.Info("Version: " + CVWONDER_VERSION)
}
