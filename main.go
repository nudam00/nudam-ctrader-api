package main

import (
	"log"
	"nudam-trading-bot/ctrader_api"
	"nudam-trading-bot/helpers/configs_helper"
)

var (
	config_path = "./configs"
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
}
