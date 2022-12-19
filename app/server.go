package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/puskipus/e-commerce/app/database/seeders"
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

func (server *Server) Initialize(appconfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appconfig.AppName)
	server.InitializeRoutes()
}

func (server *Server) InitializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed connect to database server")
	}

}

func (server *Server) dbMigrate() {
	for _, model := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

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

func (server *Server) InitCommands(config AppConfig, dbConfig DBConfig) {
	server.InitializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
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

	err := cmdApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		server.InitCommands(appconfig, dBConfig)
	} else {
		server.Initialize(appconfig, dBConfig)
		server.Run(":" + appconfig.AppPort)
	}

}
