package main

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/routes"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
)


func main(){

	storage.NewConnection()

	app:=fiber.New()

	routes.HealthRoute(app)
	routes.GetUser(app)
	routes.CreateUser(app)

	app.Listen(":3000")
}