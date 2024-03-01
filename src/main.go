package main

import (
	"log"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/runners"
)

var (
	config_path  = "./configs"
	symbolPeriod = map[string]string{"AUDCHF": "m15"}
)

func main() {
	err := configs_helper.InitializeConfig(config_path)
	if err != nil {
		log.Panic(err)
	}

	runners.TradeRoutines(symbolPeriod)
}
