package service

import (
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func ValidateUser(ctx *fiber.Ctx, email string, password string) (bool,models.User, error) {
	user, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error().Err(err).Msg("Error getting the user by email")
		return false,models.User{}, errors.New("cannot find user by email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("Incorrect Password")
		return false,models.User{}, errors.New("incorrect password")
	}

	return true,user, nil
}
