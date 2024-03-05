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

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	CTraderConfig = common.CTraderConfig{}
	if err := viper.UnmarshalKey("ctrader_config", &CTraderConfig); err != nil {
		return err
	}

	return nil
}

func initializeCTraderAccountConfig() error {
	viper.SetConfigName("ctrader_account_config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	CTraderAccountConfig = common.CTraderAccountConfig{}
	if err := viper.UnmarshalKey("ctrader_account", &CTraderAccountConfig); err != nil {
		return err
	}

	return nil
}

func initializeTraderConfiguration() error {
	viper.SetConfigName("constants")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	TraderConfiguration = common.TraderConfiguration{}
	if err := viper.UnmarshalKey("trader_configuration", &TraderConfiguration); err != nil {
		return err
	}

	return nil
}

func initializeStrategy() error {
	viper.SetConfigName("strategy")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	Strategy = common.Strategy{}
	if err := viper.UnmarshalKey("strategy", &Strategy); err != nil {
		return err
	}

	return nil
}
