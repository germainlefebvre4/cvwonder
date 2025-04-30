package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func ToggleDebug(cmd *cobra.Command, args []string) {
	if CliArgs.Debug {
		logrus.Debug("Debug logs enabled")
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		plainFormatter := new(PlainFormatter)
		logrus.SetFormatter(plainFormatter)
	}
}
