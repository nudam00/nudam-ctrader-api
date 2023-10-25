package common

type CTraderConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CTraderAccountConfig struct {
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	CtidTraderAccountId int64  `json:"ctidTraderAccountId"`
	AccessToken         string `json:"accessToken"`
}
