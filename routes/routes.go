package routes

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/controllers"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.All("/healthz/",middleware.MetricsMiddleware, controllers.CheckHealth)
	app.All("healthz/*", controllers.ErrorHealthCheck)
	app.Get("v1/user/self",middleware.MetricsMiddleware,middleware.ConnectionCheck,middleware.BasicAuthMiddleware,controllers.GetUser)
	app.Post("/v1/user",middleware.MetricsMiddleware,middleware.ConnectionCheck,controllers.CreateUser)
	app.Put("v1/user/self",middleware.MetricsMiddleware,middleware.ConnectionCheck,middleware.BasicAuthMiddleware,controllers.UpdateUser)
	app.Get("/verify",controllers.VerifyUser)
	app.Post("v1/user/self/pic",middleware.MetricsMiddleware,middleware.BasicAuthMiddleware,controllers.UploadProfilePic)
	app.Get("v1/user/self/pic",middleware.MetricsMiddleware,middleware.BasicAuthMiddleware,controllers.GetProfilePic)
	app.Delete("v1/user/self/pic",middleware.MetricsMiddleware,middleware.BasicAuthMiddleware,controllers.DeleteProfilePic)
}
