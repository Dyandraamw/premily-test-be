package app

import (
	"flag"
	_ "fmt"
	"log"
	_ "net/http"
	"os"

	"github.com/frangklynndruru/premily_backend/app/controllers"
	_ "github.com/frangklynndruru/premily_backend/app/database/seeders"
	_ "github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/urfave/cli"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

func GetEnv(key, fallback string) string   {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = controllers.Server{}
	var appConfig = controllers.AppConfig{}
	// var config =
	var dbConfig = controllers.DBConfig{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error OS loading .env file")
	}

	
	appConfig.AppName = os.Getenv("APP_NAME")
	appConfig.AppEnv = os.Getenv("APP_ENV")
	appConfig.AppPort = os.Getenv("APP_PORT")
	
	if port := os.Getenv("APP_PORT"); port != "" {
        appConfig.AppPort = port
    }

	dbConfig.DBHost = os.Getenv("DB_HOST")
	dbConfig.DBUser = os.Getenv("DB_USER")
	dbConfig.DBPassword = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	dbConfig.SSLMode = os.Getenv("DB_SSLMODE")
	
	flag.Parse()

	arg := flag.Arg(0)

	if arg != "" {
		server.InitCommands(appConfig, dbConfig)
	}else{
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}
}
