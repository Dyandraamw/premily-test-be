package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frangklynndruru/premily_backend/app/database/seeders"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("selamat datang di " + appConfig.AppName)

	// server.initializeDB(dbConfig)
	server.initializeRoutes()
	// seeders.DBSeed(server.DB)
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)

	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword,
		dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		// fmt.Printf("Gagal melakukan koneksi ke database!")
	}

	
}

func (server *Server) dbMigrate(){
	for _, model := range RegisterModels(){
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil{
			log.Fatal(err)
			fmt.Println("Gagal Migrasi!")
		}
	}

	fmt.Println("Migrasi Berhasil !")
}

//method yang digunakan untuk menjaga apabila ENV kosong
/*func getEnv(key, fallback string)string{
	if value, ok := os.LookupEnv(key); ok{
		return value
	}

	return fallback
}*/
//getEnv[namamethod]("value dari ENV", "value dari getEnv[default]")

func (server *Server) initCommands(config AppConfig, dbConfig DBConfig ){
	server. initializeDB(dbConfig) 
	commandApp	:= cli.NewApp()
	commandApp.Commands = []cli.Command{
		{
			Name : "db:migration",
			Action : func(c *cli.Context) error{
						server.dbMigrate()
						
						return nil
			},
		},
		{
			Name : "db:seed",
			Action: func(c *cli.Context)error{
					err := seeders.DBSeed(server.DB)
					if err != nil {
						log.Fatal(err)
					}
					return nil
			},
		},
	}
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error OS loading .env file")
	}

	
	appConfig.AppName = os.Getenv("APP_NAME")
	appConfig.AppEnv = os.Getenv("APP_ENV")
	appConfig.AppPort = os.Getenv("APP_PORT")
	
	dbConfig.DBHost = os.Getenv("DB_HOST")
	dbConfig.DBUser = os.Getenv("DB_USER")
	dbConfig.DBPassword = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.DBPort = os.Getenv("DB_PORT")
	
	flag.Parse()

	arg := flag.Arg(0)

	if arg != "" {
		server.initCommands()
	}else{
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}
}
