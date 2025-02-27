package config

import (
	"github.com/lenarlenar/gomart/internal/env"
	"github.com/lenarlenar/gomart/internal/flags"
)

type Config struct {
	Address        string
	DBUri          string
	AccrualAddress string
}

func Get() Config {
	envConfig, err := env.Parse()
	flagsConfig := flags.Parse()

	if(err != nil) {
		return Config{
			Address: flagsConfig.Address,
			DBUri: flagsConfig.DBUri,
			AccrualAddress: flagsConfig.AccrualAddress,
		}
	} else {
		return Config{
			Address: envConfig.Address,
			DBUri: envConfig.DBUri,
			AccrualAddress: envConfig.AccrualAddress,
		}
	}
}