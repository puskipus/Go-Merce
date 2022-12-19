package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/puskipus/e-commerce/app/database/seeders"
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

func (server *Server) Initialize(appconfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appconfig.AppName)

	server.InitializeDB(dbConfig)
	server.InitializeRoutes()
	seeders.DBSeed(server.DB)
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed connect to database server")
	}

	for _, model := range RegisterModels() {
		err = server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Migration success")
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Run() {
	var server = Server{}
	var appconfig = AppConfig{}
	var dBConfig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load env file")
	}

	appconfig.AppName = Getenv("APP_NAME", "Go-Merce")
	appconfig.AppEnv = Getenv("APP_ENV", "development")
	appconfig.AppPort = Getenv("APP_PORT", "9000")

	dBConfig.DBHost = Getenv("DB_HOST", "localhost")
	dBConfig.DBUser = Getenv("DB_USER", "user")
	dBConfig.DBPassword = Getenv("DB_PASSWORD", "password")
	dBConfig.DBName = Getenv("DB_NAME", "dbname")
	dBConfig.DBPort = Getenv("DB_PORT", "5432")

	server.Initialize(appconfig, dBConfig)
	server.Run(":" + appconfig.AppPort)
}
