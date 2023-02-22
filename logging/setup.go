package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogging(logger *logrus.Logger) {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.WarnLevel)
}
