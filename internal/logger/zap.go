package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger

func New() (*Logger, error) {
	cfg := zap.NewProductionConfig()
	var lvl zapcore.Level

	cfg.Level.SetLevel(lvl)
	cfg.DisableStacktrace = true
	cfg.Sampling.Initial = 50
	cfg.Sampling.Thereafter = 50
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.With(
		zap.String("service", "sg"),
	), nil
}
