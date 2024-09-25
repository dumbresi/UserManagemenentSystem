package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/gofiber/fiber/v2"
)

func HealthRoute(app *fiber.App){
	app.All("/healthz/",controllers.CheckHealth)
	app.All("healthz/*",controllers.ErrorHealthCheck)
}