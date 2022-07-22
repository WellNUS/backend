package config

import (
	"github.com/spf13/viper"
	"log"
)

var COOKIE_DOMAIN, FRONTEND_DOMAIN, BACKEND_DOMAIN, FRONTEND_URL, BACKEND_URL, POSTGRES_URL string
	
var (
	MATCH_THRESHOLD int = 40
	MATCH_GROUPSIZE int = 4
)

func LoadENV(path string) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil { log.Fatal(err.Error()) }

	COOKIE_DOMAIN = viper.GetString("COOKIE_DOMAIN")
	FRONTEND_DOMAIN = viper.GetString("FRONTEND_DOMAIN")
	BACKEND_DOMAIN = viper.GetString("BACKEND_DOMAIN")
	FRONTEND_URL = viper.GetString("FRONTEND_URL")
	BACKEND_URL = viper.GetString("BACKEND_URL")
	POSTGRES_URL = viper.GetString("POSTGRES_URL")
}