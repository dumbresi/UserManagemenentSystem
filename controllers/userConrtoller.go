package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/helper"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/stats"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
)

func GetUser(ctx *fiber.Ctx) error {

	if ctx.Method()==fiber.MethodPut{
		UpdateUser(ctx)
		return nil
	}

	if ctx.Method() != fiber.MethodGet {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}

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

	if ctx.Method()!=fiber.MethodPost{
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}
	var user = new(models.User)
	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	err := j.Decode(&user)
	var existingUser,er=storage.GetUserByEmail(ctx,user.Email)
	if(er ==nil){
		if(user.Email==existingUser.Email){
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error":"Email already exist"})
		}
	}
	
	if user.ID != "" || !user.AccountCreated.IsZero() || !user.AccountUpdated.IsZero() {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Request contains disallowed fields"})
    }

	if err != nil {
		log.Println("Failed to parse Response")
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Incorrect Request Body"})
	}

	if len(ctx.Queries()) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error": "Request has query parameters"})
	}

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
	startTime:= time.Now()
	err = storage.Database.Create(&user).Error
	stats.TimeDataBaseQuery("create_user",startTime,time.Now())
	if err != nil {
		log.Println("Cannot save the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user"})
	}
	fmt.Println(&user)
	return ctx.Status(http.StatusCreated).JSON(models.UserResponse{
		ID: user.ID,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		AccountCreated: user.AccountCreated,
		AccountUpdated: user.AccountUpdated,
	})
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

	if input.ID != "" || !input.AccountCreated.IsZero() || !input.AccountUpdated.IsZero() {
        return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Request contains disallowed fields"})
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
		AccountCreated: olduser.AccountCreated,
		AccountUpdated: time.Now(),
	}
	startTime:=time.Now()
	err=storage.Database.Save(&updatedUser).Error
	stats.TimeDataBaseQuery("create_user",startTime,time.Now())
	if(err!=nil){
		log.Println("Cannot update the user to Database")
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update user"})
	}
	ctx.Status(http.StatusNoContent)
	return nil
}

func ErrorPath(ctx *fiber.Ctx) error{
	ctx.Status(http.StatusNotFound)
	return nil
}
