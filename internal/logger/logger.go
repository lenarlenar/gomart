package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level, env string) error {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("ошибка парсинга уровня логирования: %w", err)
	}

	var config zap.Config
	if env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = logLevel
	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("ошибка построения логгера: %w", err)
	}

	Log = logger
	return nil
}
