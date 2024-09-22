package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

func main(){

	err:=godotenv.Load(".env")
	
	if(err!=nil){
		log.Fatal(err)
	}

	app:=fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error{
		if(c!=nil){
			return c.SendStatus(fiber.StatusOK)
		}
		return nil
	  })
	app.Listen(":3000")
}