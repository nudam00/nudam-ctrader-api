package main

import (
	"log"
	"nudam-ctrader-api/ctrader_api"
	"nudam-ctrader-api/helpers/configs_helper"
)

var (
	config_path = "./configs"
	numberDays  = 200
	period      = "D1"
	symbol      = "EURUSD"
)

func main() {
	err := configs_helper.InitializeCTraderConfig(config_path)
	if err != nil {
		log.Panic(err)
	}

	api := ctrader_api.NewCTraderAPI(numberDays, period, symbol)
	err = api.InitalizeCTrader()
	if err != nil {
		log.Panic(err)
	}

	err = api.GetTrendbars()
	if err != nil {
		log.Panic(err)
	}
}
