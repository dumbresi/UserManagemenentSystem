package awsconf

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var bucketName string

func InitS3Client() *s3.Client {

    err:=godotenv.Load("../.env")
	if(err!=nil){
		log.Error().Err(err).Msg("Error loading Env")
		return nil
	}

	bucketName= os.Getenv("S3_Bucket_Name")


	// cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("dev"),config.WithRegion("us-east-1"))
	cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithRegion(os.Getenv("AWS_Region")))
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}

func GetBucketName() string{
	log.Info().Str("Bucket Name:",bucketName)
	return bucketName
}
