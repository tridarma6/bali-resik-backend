package logger

import (
	"io"
	"os"

	"github.com/indim/bali-resik-backend/internal/config"
	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *config.LoggerConfig) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(config.GetLogLevel(cfg.Level))

	log.SetOutput(os.Stdout)

	log.SetReportCaller(true)

	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
			ForceColors:     true,
		})
	}

	return log
}

func NewNullLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(io.Discard)
	log.SetLevel(logrus.PanicLevel)
	return log
}
