package configs_helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeCTraderConfig_WrongPath(t *testing.T) {
	path := "/xyz"
	err := InitializeCTraderConfig(path)
	assert.NotNil(t, err)

	assert.Empty(t, CTraderConfig.Host)
	assert.Empty(t, CTraderConfig.Port)
	assert.Empty(t, CTraderAccountConfig.ClientId)
	assert.Empty(t, CTraderAccountConfig.ClientSecret)
	assert.Empty(t, CTraderAccountConfig.CtidTraderAccountId)
	assert.Empty(t, CTraderAccountConfig.AccessToken)
}

func TestInitializeCTraderConfig_OK(t *testing.T) {
	path := "../../configs"
	err := InitializeCTraderConfig(path)
	assert.Nil(t, err)

	assert.NotEmpty(t, CTraderConfig.Host)
	assert.NotEmpty(t, CTraderConfig.Port)
	assert.NotEmpty(t, CTraderAccountConfig.ClientId)
	assert.NotEmpty(t, CTraderAccountConfig.ClientSecret)
	assert.NotEmpty(t, CTraderAccountConfig.CtidTraderAccountId)
	assert.NotEmpty(t, CTraderAccountConfig.AccessToken)
}
