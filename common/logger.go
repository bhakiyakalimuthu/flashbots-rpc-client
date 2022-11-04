package common

import (
	"os"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	var zapCore zapcore.Core
	_cfg := zap.NewDevelopmentEncoderConfig()
	_cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCore = zapcore.NewCore(
		zapcore.NewConsoleEncoder(_cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	)
	logger := zap.New(zapCore, zap.AddCaller(), zap.ErrorOutput(zapcore.Lock(os.Stderr)))
	return logger
}
