package initserver

import (
	"os"

	"github.com/sirupsen/logrus"

	"asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/global"
)

func LoggerInit() {
	global.AsongLogger = logrus.New()

	global.AsongLogger.SetFormatter(&logrus.TextFormatter{})
	global.AsongLogger.SetOutput(os.Stdout)
	global.AsongLogger.SetLevel(logrus.DebugLevel)
}
