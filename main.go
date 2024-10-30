package main

import (
	"log"
	"os"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/middleware"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/routes"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/stats"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	storage.NewConnection()

	app := fiber.New()

	routes.SetupRoutes(app)
	
	stats.InitStatsDClient()
	app.Use(middleware.MetricsMiddleware)

	err:=godotenv.Load(".env")
	if(err!=nil){
		log.Print("Error loading Env")
		return
	}
	app.Listen(":"+os.Getenv("App_Port"))
}
