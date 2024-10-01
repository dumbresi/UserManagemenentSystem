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
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&user)

	if err != nil {
		log.Println("Failed to parse Response")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect Request Body"})
	}

	if len(ctx.Queries()) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error": "Request has query parameters"})
	}

	// validating the fields
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" {
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
	
	emailValidation:=helper.ValidateEmail(user.Email)
	if(!emailValidation){
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email not valid",
		})
	}

	user.Password = hashedPassword

	err = storage.Database.Create(&user).Error
	if err != nil {
		log.Println("Cannot save the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user"})
	}
	ctx.Status(http.StatusCreated)
	return nil
}

func UpdateUser(ctx *fiber.Ctx)error{
	var input = new(models.User)
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&input)

	if err != nil {
		log.Println("Failed to parse Response")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect Request Body"})
	}

	if len(ctx.Queries()) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error": "Request has query parameters"})
	}

	if(input.Email!=""){
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error":"Cannot change email"})
	}
	if(input.Password!=""){
		hashedPassword, err := helper.HashPassword(input.Password)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not hash password",
			})
		}
	input.Password=hashedPassword
	}
	olduser := ctx.Locals("user").(models.User)
	if(input.Password==""){
		input.Password=olduser.Password
	}
	if(input.FirstName==""){
		input.FirstName=olduser.FirstName
	}
	if(input.LastName==""){
		input.LastName=olduser.LastName
	}

	updatedUser:=models.User{
		ID: olduser.ID,
		Email: olduser.Email,
		Password: input.Password,
		FirstName: input.FirstName,
		LastName: input.LastName,
		AccountCreated: input.AccountCreated,
		AccountUpdated: input.AccountUpdated,
	}
	err=storage.Database.Save(&updatedUser).Error
	if(err!=nil){
		log.Println("Cannot update the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update user"})
	}
	return ctx.Status(http.StatusOK).JSON(models.UserResponse{
		ID: olduser.ID,
		Email: olduser.Email,
		FirstName: input.FirstName,
		LastName: input.LastName,
		AccountCreated: input.AccountCreated,
		AccountUpdated: input.AccountUpdated,
	})
	
}
