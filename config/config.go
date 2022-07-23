package config

import (
	"github.com/spf13/viper"
	"os"
	"log"
)

var COOKIE_DOMAIN, FRONTEND_DOMAIN, BACKEND_DOMAIN, FRONTEND_URL, BACKEND_URL, POSTGRES_URL string
	
var (
	MATCH_THRESHOLD int = 40
	MATCH_GROUPSIZE int = 4
)

func LoadENV(path string) {
	_, err := os.Stat(path)
	if err == nil {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err.Error())
		}
		os.Setenv("COOKIE_DOMAIN", viper.GetString("COOKIE_DOMAIN"))
		os.Setenv("FRONTEND_DOMAIN", viper.GetString("FRONTEND_DOMAIN"))
		os.Setenv("BACKEND_DOMAIN", viper.GetString("BACKEND_DOMAIN"))
		os.Setenv("FRONTEND_URL", viper.GetString("FRONTEND_URL"))
		os.Setenv("BACKEND_URL", viper.GetString("BACKEND_URL"))
		os.Setenv("POSTGRES_URL", viper.GetString("POSTGRES_URL"))
	} else {
		log.Println(err.Error())
	}

	COOKIE_DOMAIN 	= 	os.Getenv("COOKIE_DOMAIN")
	FRONTEND_DOMAIN = 	os.Getenv("FRONTEND_DOMAIN")
	BACKEND_DOMAIN 	= 	os.Getenv("BACKEND_DOMAIN")
	FRONTEND_URL 	= 	os.Getenv("FRONTEND_URL")
	BACKEND_URL 	= 	os.Getenv("BACKEND_URL")
	POSTGRES_URL 	= 	os.Getenv("POSTGRES_URL")
}