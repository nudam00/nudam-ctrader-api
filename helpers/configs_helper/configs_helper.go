package configs_helper

import (
	"log"
	"nudam-ctrader-api/types/common"

	"github.com/spf13/viper"
)

var (
	CTraderConfig        common.CTraderConfig
	CTraderAccountConfig common.CTraderAccountConfig
	TraderConfiguration  common.TraderConfiguration
	Strategy             common.Strategy
)

// Initializes cTrader config with basic variables.
func InitializeConfig(path string) error {
	log.Printf("initializes config...")

	viper.Reset()
	viper.AddConfigPath(path)
	viper.SetConfigType("json")

	if err := initializeCTraderConfig(); err != nil {
		return err
	}

	if err := initializeCTraderAccountConfig(); err != nil {
		return err
	}

	if err := initializeTraderConfiguration(); err != nil {
		return err
	}

	if err := initializeStrategy(); err != nil {
		return err
	}

	return nil
}

func initializeCTraderConfig() error {
	viper.SetConfigName("ctrader_demo_config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	CTraderConfig = common.CTraderConfig{}
	err = viper.UnmarshalKey("ctrader_config", &CTraderConfig)
	if err != nil {
		return err
	}

	return nil
}

func initializeCTraderAccountConfig() error {
	viper.SetConfigName("ctrader_account_config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	CTraderAccountConfig = common.CTraderAccountConfig{}
	err = viper.UnmarshalKey("ctrader_account", &CTraderAccountConfig)
	if err != nil {
		return err
	}

	return nil
}

func initializeTraderConfiguration() error {
	viper.SetConfigName("constants")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	TraderConfiguration = common.TraderConfiguration{}
	err = viper.UnmarshalKey("trader_configuration", &TraderConfiguration)
	if err != nil {
		return err
	}

	return nil
}

func initializeStrategy() error {
	viper.SetConfigName("strategy")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	Strategy = common.Strategy{}
	err = viper.UnmarshalKey("strategy", &Strategy)
	if err != nil {
		return err
	}

	return nil
}
