package controllers

import (
	"log"
	"net/http"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
)


func GetUser(ctx *fiber.Ctx)error{
	return nil
}
func CreateUser(ctx *fiber.Ctx)error{
	var user= new(models.User)

	err := ctx.BodyParser(user)
	if err !=nil{
		log.Fatal("Failed to parse Response")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"Incorrect Request Body"})
	}
	err=storage.Database.Create(&user).Error
	if err!=nil{
		log.Fatal("Cannot save the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"Cannot create user"})
	}
	ctx.Status(http.StatusOK)
	return nil
}