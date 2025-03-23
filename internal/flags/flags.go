package flags

import (
	"flag"
)

const (
	DefaultAddress        = ":8081"
	DefaultDBUri          = "host=localhost port=5432 user=gomart password=gomart dbname=gomart sslmode=disable"
	DefaultAccrualAddress = "http://localhost:8080"
)

type FlagsConfig struct {
	Address        string
	DBUri          string
	AccrualAddress string
}

func Parse() FlagsConfig {
	var c FlagsConfig
	flag.StringVar(&c.Address, "a", DefaultAddress, "адрес и порт запуска")
	flag.StringVar(&c.DBUri, "d", DefaultDBUri, "адрес подключения к БД")
	flag.StringVar(&c.AccrualAddress, "r", DefaultAccrualAddress, "адрес системы расчёта начислений")
	flag.Parse()
	return c
}
