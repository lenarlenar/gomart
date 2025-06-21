package env

import (
	"fmt"

	"github.com/caarlos0/env"
)

type EnvConfig struct {
	Address        string `env:"RUN_ADDRESS"`
	DBUri          string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Parse() (EnvConfig, error) {
	var envConfig EnvConfig
	if err := env.Parse(&envConfig); err != nil {
		return EnvConfig{}, fmt.Errorf("ошибка при парсинге переменные окружения: %v", err)
	}

	return envConfig, nil
}
