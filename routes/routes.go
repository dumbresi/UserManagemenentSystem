package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/gofiber/fiber/v2"
)

func HealthRoute(app *fiber.App){
	app.Get("/healthz",controllers.CheckHealth)
}