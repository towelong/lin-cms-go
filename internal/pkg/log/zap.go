package log

import (
	"path"

	"github.com/spf13/viper"
	"github.com/towelong/lin-cms-go/pkg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func NewCustomerLogger() {
	stdout := []string{"stdout"}
	stderr := []string{"stderr"}
	env := viper.GetString("env")
	if env == "prod" {
		logPath, _ := pkg.CreateDirAndFileForCurrentTime("logs", "2006-01-02")
		stdLog := path.Join(logPath, "/log.txt")
		errorLog := path.Join(logPath, "/errLog.txt")
		stdout = append(stdout, stdLog)
		stderr = append(stderr, errorLog)
	}
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: env == "dev",
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
		},
		OutputPaths:      stdout,
		ErrorOutputPaths: stderr,
	}
	Logger, _ = config.Build()
	defer Logger.Sync()
}
