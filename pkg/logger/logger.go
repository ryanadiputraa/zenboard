package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLoger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
