package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ErrorLogger() KafkaLogger {
	return KafkaLogger{logger: log.Logger}
}

type KafkaLogger struct {
	logger zerolog.Logger
}

func (l KafkaLogger) Printf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}
