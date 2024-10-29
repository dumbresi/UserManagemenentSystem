package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

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

	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()
	j.Decode(&models.Image{})

	file, err := ctx.FormFile("profilePic")
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).SendString("No file uploaded")
    }
	existingImage, er:=storage.GetProfilePicByUserId(ctx,user.ID)
	
	if(er==nil && existingImage.UserID!=""){
		storage.DeleteProfilePicbyId(ctx,existingImage.ID)
		DeleteExistingPic(s3Client,bucketName, existingImage)
	}

	s3URL, err := UploadToS3(s3Client, file,user.ID)

	if err != nil {
		log.Println("Failed to upload image data in bucket")
        return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to upload image")
    }

	var image models.Image

	image.UserID= user.ID
	image.URL=s3URL
	image.FileName=file.Filename

    if err := storage.Database.Create(&image).Error; err != nil {
		log.Println("Failed to save image data in DB")
        return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save image data")
    }
	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"file_name": image.FileName,
		"id": image.ID,
		"url": image.URL,
		"upload_date": image.UploadDate,
		"user_id": image.UserID,
	})
	return nil
}

func GetProfilePic(ctx *fiber.Ctx) error{
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

	
	var user= ctx.Locals("user").(models.User)
	profilePic,err:= storage.GetProfilePicByUserId(ctx,user.ID)
	if(err!=nil){
		ctx.Status(http.StatusNotFound)
		return nil
	}

	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"file_name": profilePic.FileName,
		"id": profilePic.ID,
		"url": profilePic.URL,
		"upload_date": profilePic.UploadDate,
		"user_id": profilePic.UserID,
	})

	return nil
}

func DeleteProfilePic(ctx *fiber.Ctx) error{
	if ctx.Method() != fiber.MethodDelete {
		ctx.Status(fiber.StatusMethodNotAllowed)
		return nil
	}

	if len(ctx.Queries()) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Error": "Request has query parameters"})
	}

	if len(ctx.Body())>0 {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"Bad Request with error" : "Request has a payload"})
	}

	var user= ctx.Locals("user").(models.User)
	profilePic,err:= storage.GetProfilePicByUserId(ctx,user.ID)
	if(err!=nil){
		ctx.Status(http.StatusNotFound)
		return nil
	}
	err=DeleteExistingPic(s3Client,bucketName,profilePic)
	if(err!=nil){
		log.Print("Cannot delete profilePic from Bucket")
		ctx.Status(http.StatusInternalServerError)
		return nil
	}

	err=storage.DeleteProfilePicbyId(ctx,profilePic.ID)
	if(err!=nil){
		log.Println("Cannot delete profilePic from DB")
		ctx.Status(http.StatusInternalServerError)
		return nil
	}
	ctx.Status(http.StatusNoContent)
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

func DeleteExistingPic(s3Client *s3.Client,bucketName string, image models.Image) error{
	objectKey:=image.UserID+"/"+image.ID;
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: &bucketName,
        Key:    &objectKey,
    })
    if err != nil {
        return err
    }
	return nil
}