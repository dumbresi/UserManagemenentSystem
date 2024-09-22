package main

import (
	"log"
	"os"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

func(r *Repository) SetUpRoutes(app *fiber.App){
	app.Get("/healthz", func(c *fiber.Ctx) error{
		c.Status(fiber.StatusOK)
		return nil
	})
}
func main(){

	err:=godotenv.Load(".env")
	
	if(err!=nil){
		log.Fatal(err)
	}

	config:=&storage.Config{
		Host: os.Getenv("DB_Host"),
		Port: os.Getenv("DB_Port"),
		User: os.Getenv("DB_User"),
		Password: os.Getenv("DB_Password"),
		DbName: os.Getenv("DB_Name"),
		SSLMode: os.Getenv("DB_SslMode"),
	}
	db,err:= storage.NewConnection(config)

	if(err!=nil){
		log.Fatal(err.Error())
	}
	r:=Repository{
		DB: db,
	}
	app:=fiber.New()
	r.SetUpRoutes(app)
	app.Listen(":3000")
}