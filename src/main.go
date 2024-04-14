package main

import (
	"log"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/runners"
)

var (
	config_path = "./configs"
)

func main() {
	if err := configs_helper.InitializeConfig(config_path); err != nil {
		log.Panic(err)
	}
	runners.TradeRoutines()
}
