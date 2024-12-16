package utils

import (
	"github.com/sirupsen/logrus"
)

func CheckError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
