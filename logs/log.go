package logs

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logConfig.CallerKey = "caller"
	logConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.InfoLevel,
	)
	logger = zap.New(core)

}

func I() *zap.Logger {
	return logger
}
