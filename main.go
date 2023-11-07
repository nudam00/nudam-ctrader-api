package main

import (
	"log"
	"nudam-ctrader-api/ctrader_api"
	"nudam-ctrader-api/helpers/configs_helper"
	"nudam-ctrader-api/strategy"
)

var (
	config_path = "./configs"
	period      = "m15"
	symbol      = "EURUSD"
)

func main() {
	err := configs_helper.InitializeConfig(config_path)
	if err != nil {
		log.Panic(err)
	}

	api := ctrader_api.NewCTraderAPI(period, symbol)
	err = api.InitalizeCTrader()
	if err != nil {
		log.Panic(err)
	}

	trader := strategy.NewTrader()
	err = api.GetTrendbars(trader)
	if err != nil {
		log.Panic(err)
	}
}
