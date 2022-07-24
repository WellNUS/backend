package config

import (
	"github.com/spf13/viper"
	"os"
	"log"
)

var COOKIE_ADDRESS, SERVER_ADDRESS, FRONTEND_ADDRESS, BACKEND_ADDRESS, DB_ADDRESS string
	
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
		os.Setenv("FRONTEND_ADDRESS", 	viper.GetString("FRONTEND_ADDRESS"))
		os.Setenv("BACKEND_ADDRESS", 	viper.GetString("BACKEND_ADDRESS"))

		os.Setenv("DB_ADDRESS", 		viper.GetString("DB_ADDRESS"))
		os.Setenv("SERVER_ADDRESS", 	viper.GetString("SERVER_ADDRESS"))
		os.Setenv("COOKIE_ADDRESS", 	viper.GetString("COOKIE_ADDRESS"))
	} else {
		log.Println(err.Error())
	}

	FRONTEND_ADDRESS	= 	os.Getenv("FRONTEND_ADDRESS")
	BACKEND_ADDRESS 	= 	os.Getenv("BACKEND_ADDRESS")
	
	DB_ADDRESS 			= 	os.Getenv("DB_ADDRESS")
	SERVER_ADDRESS		= 	os.Getenv("SERVER_ADDRESS")
	COOKIE_ADDRESS 		= 	os.Getenv("COOKIE_ADDRESS")

	// FOR HEROKU ONLY
	port, ok := os.LookupEnv("PORT")
    if ok {
        SERVER_ADDRESS = ":" + port
    }
}