package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCustomerLogger() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	log, _ := config.Build()
	return log
}
