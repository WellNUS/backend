package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var COOKIE_ADDRESS, SERVER_ADDRESS, FRONTEND_ADDRESS, BACKEND_ADDRESS, WS_ADDRESS, DB_ADDRESS string

var (
	MATCH_THRESHOLD int = 40
	MATCH_GROUPSIZE int = 4
)

// With docker compose
func LoadENV(path string) {
	_, err := os.Stat(path)
 	if err == nil {
 		viper.SetConfigFile(path)
 		if err := viper.ReadInConfig(); err != nil {
 			log.Fatal(err.Error())
 		}

		if _, ok := os.LookupEnv("FRONTEND_ADDRESS"); !ok {
			log.Println("NOTE: Using FRONTEND_ADDRESS from .env")
			os.Setenv("FRONTEND_ADDRESS", viper.GetString("FRONTEND_ADDRESS"))
		}
		if _, ok := os.LookupEnv("BACKEND_ADDRESS"); !ok {
			log.Println("NOTE: Using BACKEND_ADDRESS from .env")
			os.Setenv("BACKEND_ADDRESS", viper.GetString("BACKEND_ADDRESS"))
		}
		if _, ok := os.LookupEnv("WS_ADDRESS"); !ok {
			log.Println("NOTE: Using WS_ADDRESS from .env")
			os.Setenv("WS_ADDRESS",	viper.GetString("WS_ADDRESS"))
		}
		if _, ok := os.LookupEnv("SERVER_ADDRESS"); !ok {
			log.Println("NOTE: Using SERVER_ADDRESS from .env")
			os.Setenv("SERVER_ADDRESS", viper.GetString("SERVER_ADDRESS"))
		}
		if _, ok := os.LookupEnv("COOKIE_ADDRESS"); !ok {
			log.Println("NOTE: Using COOKIE_ADDRESS from .env")
			os.Setenv("COOKIE_ADDRESS", viper.GetString("COOKIE_ADDRESS"))
		}
		if _, ok := os.LookupEnv("DB_ADDRESS"); !ok {
			log.Println("NOTE: Using DB_ADDRESS from .env")
			os.Setenv("DB_ADDRESS", viper.GetString("DB_ADDRESS"))
		}

 	} else {
 		log.Println(err.Error())
 	}

	FRONTEND_ADDRESS = os.Getenv("FRONTEND_ADDRESS")
	BACKEND_ADDRESS = os.Getenv("BACKEND_ADDRESS")
	WS_ADDRESS = os.Getenv("WS_ADDRESS")
	SERVER_ADDRESS = os.Getenv("SERVER_ADDRESS")
	COOKIE_ADDRESS = os.Getenv("COOKIE_ADDRESS")
	DB_ADDRESS = os.Getenv("DB_ADDRESS")

	// FOR HEROKU ONLY
	port, ok := os.LookupEnv("PORT")
	if ok {
		SERVER_ADDRESS = ":" + port
	}
}