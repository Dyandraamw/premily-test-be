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

	
	appConfig.AppName = GetEnv("APP_NAME", "Premily")
	appConfig.AppEnv = GetEnv("APP_ENV", "Development")
	appConfig.AppPort = GetEnv("APP_PORT", "9000")
	
	dbConfig.DBHost = GetEnv("DB_HOST", "localhost")
	dbConfig.DBUser = GetEnv("DB_USER", "user")
	dbConfig.DBPassword = GetEnv("DB_PASSWORD", "post9090")
	dbConfig.DBName = GetEnv("DB_NAME", "premily")
	dbConfig.DBPort = GetEnv("DB_PORT", "5433")
	
	flag.Parse()

	arg := flag.Arg(0)

	if arg != "" {
		server.InitCommands(appConfig, dbConfig)
	}else{
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}
}
