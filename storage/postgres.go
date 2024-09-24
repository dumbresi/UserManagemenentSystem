package storage

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct{
	Host string
	Port string
	User string
	Password string
	DbName string
	SSLMode string
}

var db *gorm.DB

func NewConnection() {

	err:=godotenv.Load(".env")
	if(err!=nil){
		log.Fatal()
	}

	config:=Config{
		Host: os.Getenv("DB_Host"),
		Port: os.Getenv("DB_Port"),
		User: os.Getenv("DB_User"),
		Password: os.Getenv("DB_Password"),
		DbName: os.Getenv("DB_Name"),
		SSLMode: os.Getenv("DB_SslMode"),
	}

	dsn:=fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,config.Port,config.User,config.Password,config.DbName,config.SSLMode,
	)
	var er error
	db, er=gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if er != nil {
        log.Fatal("Failed to connect to the database", err)
    }

}