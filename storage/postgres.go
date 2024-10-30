package storage

import (
	"errors"
	"fmt"
	"os"
	"github.com/rs/zerolog/log"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

var Database *gorm.DB

func NewConnection() error {

	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Err(err).Msg("error loading env")
		return err
	}

	config := Config{
		Host:     os.Getenv("DB_Host"),
		Port:     os.Getenv("DB_Port"),
		User:     os.Getenv("DB_User"),
		Password: os.Getenv("DB_Password"),
		DbName:   os.Getenv("DB_Name"),
		SSLMode:  os.Getenv("DB_SslMode"),
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DbName, config.SSLMode,
	)
	var er error
	Database, er = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if er != nil {
		log.Print("Failed to connect to the database", err)
		log.Error().Err(err).Msg("Failed to connect to the database")
		return er
	}
	er=MigrateDb()
	if(er!=nil){
		log.Error().Err(err).Msg("Error migrating the database")
	}

	return nil
}

func MigrateDb() error {
	err:=models.MigrateUser(Database)
	if(err!=nil){
		log.Print("Cannot migrate User DB")
		log.Error().Err(err).Msg("Cannt migrate Users table")
    }
	err= models.MigrateImage(Database)
	if(err!=nil){
		log.Error().Err(err).Msg("Cannt migrate Images table")
	}
	return err
}

func PingDb() error {
	sqlDB, err := Database.DB()
	if err != nil {
		log.Error().Err(err).Msg("failed to retreive database")
		return fmt.Errorf("failed to retrieve database object: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Error().Err(err).Msg("failed to ping database")
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func GetUserByEmail(ctx *fiber.Ctx, email string) (models.User, error) {
	if Database == nil {
		log.Error().Msg("DB object is not initialized")
		return models.User{}, errors.New("DB object is not initialized")
	}

	var user models.User
	err := Database.Where("email = ?", email).First(&user).Error; 
	if err != nil {
		log.Print("Error getting the user by username")
		return models.User{}, err
	}
	return user, nil
}

func DeleteProfilePicbyId(ctx *fiber.Ctx, imageId string) error{
	if Database == nil {
		log.Error().Msg("DB object is not initialized")
		return errors.New("DB object is not initialized")
	}
	var image models.Image
	err:= Database.Where("id = ?",imageId).Delete(&image).Error
	if(err!=nil){
		return err
	}

	return nil
}

func GetProfilePicByUserId(ctx *fiber.Ctx, userId string) (models.Image, error){
	if Database == nil {
		log.Error().Msg("DB object is not initialized")
		return models.Image{}, errors.New("DB object is not initialized")
	}

	var image models.Image
	err:= Database.Where("user_id = ?",userId).First(&image).Error
	if(err!=nil){
		log.Info().Msg("Pofile pic does not exist for this user")
		return image, err
	}
	
	return image, nil
}