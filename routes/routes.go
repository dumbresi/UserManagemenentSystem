package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.All("/healthz/", controllers.CheckHealth)
	app.All("healthz/*", controllers.ErrorHealthCheck)
	app.All("v1/user/self",middleware.ConnectionCheck,middleware.BasicAuthMiddleware,controllers.GetUser)
	app.All("/v1/user",middleware.ConnectionCheck,controllers.CreateUser)
	app.Post("v1/user/self/pic",middleware.BasicAuthMiddleware,controllers.UploadProfilePic)
	app.Get("v1/user/self/pic",middleware.BasicAuthMiddleware,controllers.GetProfilePic)
	app.Delete("v1/user/self/pic",middleware.BasicAuthMiddleware,controllers.DeleteProfilePic)
}
