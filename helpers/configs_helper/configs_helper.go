package configs_helper

import (
	"log"
	"nudam-trading-bot/types/common"

	"github.com/spf13/viper"
)

var (
	CTraderConfig        common.CTraderConfig
	CTraderAccountConfig common.CTraderAccountConfig
)

func InitializeCTraderConfig(path string) error {
	log.Printf("initializes config...")

	viper.Reset()

	CTraderConfig = common.CTraderConfig{}

	viper.AddConfigPath(path)
	viper.SetConfigName("ctrader_demo_config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("ctrader_config", &CTraderConfig)
	if err != nil {
		return err
	}

	CTraderAccountConfig = common.CTraderAccountConfig{}

	viper.SetConfigName("ctrader_account_config")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("ctrader_account", &CTraderAccountConfig)
	if err != nil {
		return err
	}

	return nil
}
