package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/helper"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
)

func GetUser(ctx *fiber.Ctx) error {
	
	if len(ctx.Body())>0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Bad Request with error" : "Request has a payload"})
	}

	if len(ctx.Queries()) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error": "Request has query parameters"})
	}
	user := ctx.Locals("user").(models.User)
	userResp:=models.UserResponse{
		ID: user.ID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		AccountCreated: user.AccountCreated,
		AccountUpdated: user.AccountUpdated,
	}
	return ctx.Status(http.StatusOK).JSON(userResp)
}
func CreateUser(ctx *fiber.Ctx) error {
	var user = new(models.User)

	// err := ctx.BodyParser(user)
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&user)

	if err != nil {
		log.Fatal("Failed to parse Response")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect Request Body"})
	}

	// validating the fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
		log.Fatal("firstname:"+user.FirstName)
		log.Fatal("lastname:"+user.LastName)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First name, last name, email, and password are required fields",
		})
	}
	hashedPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not hash password",
		})
	}

	user.Password = hashedPassword

	err = storage.Database.Create(&user).Error
	if err != nil {
		log.Fatal("Cannot save the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user"})
	}
	ctx.Status(http.StatusCreated)
	return nil
}
