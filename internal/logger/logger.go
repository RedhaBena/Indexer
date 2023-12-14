package logger

import (
	"context"
	"fmt"
	"time"

	"RedhaBena/indexer/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var GlobalLogger zap.Logger

func InitGlobalLogger(context context.Context) error {
	var lvl zapcore.Level
	if err := lvl.Set("debug"); err != nil {
		return fmt.Errorf("invalid log-level %q: %v", config.GlobalConfig.LoggerConfig.LogLevel, err)
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("Failed to create logger: %v\n", err)
	}
	GlobalLogger = *logger
	defer GlobalLogger.Sync()

	return nil
}
