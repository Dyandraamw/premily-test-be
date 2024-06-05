package controllers

import (
	"fmt"
	"log"
	_ "sync"

	"net/http"
	"os"

	"github.com/frangklynndruru/premily_backend/app/database/seeders"
	"github.com/frangklynndruru/premily_backend/app/models"

	_ "golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/cli"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB        *gorm.DB
	Router    *mux.Router
	AppConfig *AppConfig
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
	AppURL  string
}

type IDGenerator struct {
	prefix string
	count  int
	// mu 	sync.Mutex
}
type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var sessionUser = "user-session"

func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("selamat datang di " + appConfig.AppName)

	server.initializeDB(dbConfig)
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

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
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

func (server *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
	server.initializeDB(dbConfig)
	commandApp := cli.NewApp()
	commandApp.Commands = []cli.Command{
		{
			Name: "db:migration",
			Action: func(c *cli.Context) error {
				server.dbMigrate()

				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}
	err := commandApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Melakukan pengecekan apakah user sudah sukses untuk login atau tidak
=========================
*/
func IsLogin(r *http.Request) bool {
	session, _ := store.Get(r, sessionUser)
	if session.Values["user_id"] == nil {
		return false
	}
	return true
}

func (server *Server) CurrentUser(w http.ResponseWriter, r *http.Request) *models.User {

	if !IsLogin(r) {
		return nil
	}

	session, _ := store.Get(r, sessionUser)

	userModel := models.User{}
	user, err := userModel.FindByID(server.DB, session.Values["user_id"].(string))
	if err != nil {
		session.Values["user_id"] = nil
		session.Save(r, w)
		return nil
	}

	return user
}

// func MakePassword(password string)(string, error){
// 	comparedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(comparedPassword), err
// }

// func verifyPassword(password string, comparePassword string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(comparePassword), []byte(password)) == nil
// }

func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix: prefix,
		count:  0,
	}
}

func (g *IDGenerator) NextID() string {
	g.count = g.count + 1
	return fmt.Sprintf("%s-%s", g.prefix, g.pad(g.count, 5))
}

// func NewIDGenerator(prefix string) *IDGenerator {
// 	return &IDGenerator{
// 		prefix:  prefix,
// 		count: 0,
// 	}
// }

// func (g *IDGenerator) NextID() string {
// 	g.mu.Lock()
// 	defer g.mu.Unlock()

//		g.count++
//		return fmt.Sprintf("%s-%05d", g.prefix, g.count)
//	}
func (g *IDGenerator) pad(number, width int) string {
	return fmt.Sprintf("%0*d", width, number)
}
