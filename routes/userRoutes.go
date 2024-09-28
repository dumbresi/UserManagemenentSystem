package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/gofiber/fiber/v2"
)

func GetUser(app *fiber.App){
	app.Get("v1/user/self",controllers.GetUser)
}

func CreateUser(app *fiber.App){
	app.Post("/v1/user",controllers.CreateUser)
}