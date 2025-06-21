package config

import (
	"github.com/lenarlenar/gomart/internal/env"
	"github.com/lenarlenar/gomart/internal/flags"
)

type Config struct {
	Address        string
	DSN            string
	AccrualAddress string
}

func Get() Config {
	envConfig, err := env.Parse()
	flagsConfig := flags.Parse()

	result := Config{
		Address:        flagsConfig.Address,
		DSN:            flagsConfig.DBUri,
		AccrualAddress: flagsConfig.AccrualAddress,
	}

	if err != nil {
		return result
	}

	if envConfig.AccrualAddress != "" {
		result.AccrualAddress = envConfig.AccrualAddress
	} else {
		result.AccrualAddress = flagsConfig.AccrualAddress
	}

	if envConfig.Address != "" {
		result.Address = envConfig.Address
	} else {
		result.Address = flagsConfig.Address
	}

	if envConfig.DBUri != "" {
		result.DSN = envConfig.DBUri
	} else {
		result.DSN = flagsConfig.DBUri
	}

	return result
}
