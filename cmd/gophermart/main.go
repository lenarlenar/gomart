package main

import (
	"github.com/lenarlenar/gomart/internal"
	"github.com/lenarlenar/gomart/internal/config"
	"github.com/lenarlenar/gomart/internal/router"
)

func main() {
	internal.InitLogger()
	config := config.Get()
	router.Run(config.Address)
}
