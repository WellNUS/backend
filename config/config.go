package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var COOKIE_ADDRESS, SERVER_ADDRESS, FRONTEND_ADDRESS, BACKEND_ADDRESS, WS_ADDRESS, DB_ADDRESS string

var RUN_WITH_DOCKER_COMPOSE bool
var DOCKER_COMPOSE_DB_ADDRESS string


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
		if !viper.GetBool("IGNORE") {
			log.Println("NOTE: Writing environment variable with .env")
			os.Setenv("FRONTEND_ADDRESS", viper.GetString("FRONTEND_ADDRESS"))
			os.Setenv("BACKEND_ADDRESS", viper.GetString("BACKEND_ADDRESS"))
			os.Setenv("WS_ADDRESS",	viper.GetString("WS_ADDRESS"))
			os.Setenv("SERVER_ADDRESS", viper.GetString("SERVER_ADDRESS"))
			os.Setenv("COOKIE_ADDRESS", viper.GetString("COOKIE_ADDRESS"))

			os.Setenv("DB_ADDRESS", viper.GetString("DB_ADDRESS"))
			os.Setenv("RUN_WITH_DOCKER_COMPOSE", viper.GetString("RUN_WITH_DOCKER_COMPOSE"))
			os.Setenv("DOCKER_COMPOSE_DB_ADDRESS", viper.GetString("DOCKER_COMPOSE_DB_ADDRESS"))
		} else {
			log.Println("NOTE: .env is ignored. No writing of environment variables with .env")
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
	RUN_WITH_DOCKER_COMPOSE = os.Getenv("RUN_WITH_DOCKER_COMPOSE") == "TRUE"
	DOCKER_COMPOSE_DB_ADDRESS = os.Getenv("DOCKER_COMPOSE_DB_ADDRESS")

	// FOR HEROKU ONLY
	port, ok := os.LookupEnv("PORT")
	if ok {
		SERVER_ADDRESS = ":" + port
	}
}