package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/service"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
)

func BasicAuthMiddleware(ctx *fiber.Ctx) error{
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header missing",
		})
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization method",
		})
	}
	
	base64Credentials := strings.TrimPrefix(authHeader, "Basic ")

	decodedCredentials, err := base64.StdEncoding.DecodeString(base64Credentials)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid base64 encoding",
		})
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format",
		})
	}

	username := credentials[0]
	password := credentials[1]
	exist,user,validationerror:=service.ValidateUser(ctx,username,password); 
	if exist {
		if(!user.IsVerified){
			return ctx.SendStatus(http.StatusForbidden)
		}
		ctx.Locals("user",user)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": validationerror,
	})


}

func ConnectionCheck(ctx *fiber.Ctx)error{
	err:=storage.PingDb()
	if(err!=nil){
		ctx.Status(fiber.StatusServiceUnavailable)
		return nil
	}
	return ctx.Next()
}