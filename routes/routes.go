package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.All("/healthz/", controllers.CheckHealth)
	app.All("healthz/*", controllers.ErrorHealthCheck)
	app.Get("v1/user/self",middleware.BasicAuthMiddleware,controllers.GetUser)
	app.Post("/v1/user",controllers.CreateUser)
	app.Put("v1/user/self",middleware.BasicAuthMiddleware,controllers.UpdateUser)
}
