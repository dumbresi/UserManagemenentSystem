package awsconf

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func InitS3Client() *s3.Client {
    
	cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("dev"),config.WithRegion("us-east-1"))
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}
