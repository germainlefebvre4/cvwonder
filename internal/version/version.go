package version

import (
	"github.com/sirupsen/logrus"
)

var CVWONDER_VERSION = "0.3.1"

func (t *VersionService) GetVersion() {
	logrus.Debug("GetVersion")
	logrus.Info("Version: " + CVWONDER_VERSION)
}
