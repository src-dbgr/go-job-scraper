package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(level string) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(logLevel)
}

func LogInfo(message string) {
	log.Info().Msg(message)
}
