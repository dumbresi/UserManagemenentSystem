package controllers

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/awsconf"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

var s3Client= awsconf.InitS3Client()
var bucketName string= "sidd1234"

func UploadProfilePic(ctx *fiber.Ctx) error{
	user, ok := ctx.Locals("user").(models.User)
	if !ok {
        return ctx.Status(fiber.StatusUnauthorized).SendString("User not found")
    }

	file, err := ctx.FormFile("profilePic")
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).SendString("No file uploaded")
    }

	s3URL, err := UploadToS3(s3Client, file,user.ID)
	if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to upload image to S3")
    }
	var image models.Image

	image.UserID= user.ID
	image.URL=s3URL
	image.FileName=file.Filename

	
    if err := storage.Database.Create(&image).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save image data in database")
    }

	return nil
}


func UploadToS3(s3Client *s3.Client,file *multipart.FileHeader, userid string) (string, error){
	fileContent, err:=file.Open()
	if err != nil {
		log.Print(err.Error())
        return "", err
    }
	defer fileContent.Close()
	key:= fmt.Sprintf("%s/%s",userid,file.Filename)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
        Body:   fileContent,
	}) 

	if err != nil {
		log.Print(err.Error())
        return "", err
    }

    s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
    return s3URL, nil
}