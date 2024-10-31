package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/awsconf"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/helper"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/models"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/stats"
	"github.com/CSYE-6225-CLOUD-SIDDHARTH/webapp/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

var s3Client= awsconf.InitS3Client()
var bucketName string= awsconf.GetBucketName()

func UploadProfilePic(ctx *fiber.Ctx) error{
	user, ok := ctx.Locals("user").(models.User)
	if !ok {
        return ctx.Status(fiber.StatusUnauthorized).SendString("User not found")
    }

	j := json.NewDecoder(strings.NewReader(string(ctx.Body())))
	j.DisallowUnknownFields()

	var image models.Image

	form, err := ctx.MultipartForm()
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
    }

    files := form.File["profilePic"]
    if len(files) == 0 {
        return ctx.Status(fiber.StatusBadRequest).SendString("No file uploaded")
    }
    if len(files) > 1 {
        return ctx.Status(fiber.StatusBadRequest).SendString("Only one image file is allowed")
    }

    file := files[0]

    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).SendString("No file uploaded")
    }
	existingImage, er:=storage.GetProfilePicByUserId(ctx,user.ID)
	
	if(er==nil && existingImage.UserID!=""){
		return ctx.Status(fiber.StatusConflict).SendString("Profile Picture already exists")
	}

	err=helper.ValidateImageFile(file)
	if(err!=nil){
		return ctx.Status(fiber.StatusBadRequest).SendString("Please upload an Image file format")
	}

	s3URL, err := UploadToS3(s3Client, file,user.ID)

	if err != nil {
		log.Error().Msg("Failed to upload image data in bucket")
        return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to upload image")
    }

	image.UserID= user.ID
	image.URL=s3URL
	image.FileName=file.Filename
	startTime:= time.Now()
    err = storage.Database.Create(&image).Error; 
	stats.TimeDataBaseQuery("create_pic",startTime,time.Now())
	if err != nil {
		log.Error().Msg("Failed to save image data in DB")
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
		log.Error().Err(err).Msg("User not found")
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
	log.Print("Before deleting existing pic")
	err=DeleteExistingPic(ctx,s3Client,bucketName,profilePic)
	if(err!=nil){
		log.Print("Cannot delete profilePic from Bucket")
		ctx.Status(http.StatusInternalServerError)
		return nil
	}

	err=storage.DeleteProfilePicbyId(ctx,profilePic.ID)
	if(err!=nil){
		log.Error().Msg("Cannot delete profilePic from DB")
		ctx.Status(http.StatusInternalServerError)
		return nil
	}
	ctx.Status(http.StatusNoContent)
	return nil
}

func UploadToS3(s3Client *s3.Client,file *multipart.FileHeader, userid string) (string, error){
	fileContent, err:=file.Open()
	if err != nil {
		log.Error().Err(err).Msg("Error opening the image file")
        return "", err
    }
	defer fileContent.Close()
	key:= fmt.Sprintf("%s/%s",userid,file.Filename)
	startTime:=time.Now()
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
        Body:   fileContent,
	}) 
	stats.TimeS3Call("put_image",startTime,time.Now())

	if err != nil {
		log.Print(err.Error())
        return "", err
    }

	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s3Client.Options().Region, key)
    return s3URL, nil
}

func DeleteExistingPic(ctx *fiber.Ctx,s3Client *s3.Client,bucketName string, image models.Image) error{
	log.Printf("delete Existing Pic %s : %v",bucketName,s3Client)
	objectKey:=fmt.Sprintf(image.UserID+"/"+image.FileName);
	startTime:=time.Now()
	resp, err := s3Client.DeleteObject(ctx.Context(), &s3.DeleteObjectInput{
        Bucket: &bucketName,
        Key:    &objectKey,
    })
	log.Printf("Response of delete %v",resp.ResultMetadata)
	log.Printf("delete object %s : %s",bucketName,objectKey)
	stats.TimeS3Call("delete_image",startTime,time.Now())
    if err != nil {
		log.Printf("Not able to delete object %s : %s",bucketName,objectKey)
		log.Error().Str("Objkey",objectKey).Msg("Error deleting the pic from bucket")
        return err
    }
	return nil
}