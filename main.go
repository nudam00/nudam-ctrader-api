package main

import (
	"log"
	"nudam-ctrader-api/ctrader_api"
	"nudam-ctrader-api/helpers/configs_helper"
)

var (
	config_path = "./configs"
	numberDays  = 30
	period      = "M15"
	symbol      = "EURUSD"
	countBars   = uint32(50)
)

func main() {
	err := configs_helper.InitializeCTraderConfig(config_path)
	if err != nil {
		log.Panic(err)
	}

	api := ctrader_api.NewCTraderAPI()
	err = api.InitializeWsDialer()
	if err != nil {
		log.Panic(err)
	}

	err = api.SaveAvailableSymbols()
	if err != nil {
		log.Panic(err)
	}

	err = api.GetTrendbars(numberDays, period, symbol, countBars)
	if err != nil {
		log.Panic(err)
	}
}
