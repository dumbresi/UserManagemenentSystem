package awsconf

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var bucketName string
var s3Client *s3.Client
var snsTopicArn string
var snsClient *sns.Client

func InitS3Client() *s3.Client {

    err:=godotenv.Load(".env")
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
	if(s3Client==nil){
		s3Client=s3.NewFromConfig(cfg)
	}
	return s3Client
}

func GetS3Client() *s3.Client{
	return s3Client
}

func GetBucketName() string{
	log.Info().Str("Bucket Name:",bucketName)
	return bucketName
}



func InitSNSClient(){
	// cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithSharedConfigProfile("dev"),config.WithRegion("us-east-1"))
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_Region")))

	if err != nil {
		log.Error().Msg(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	snsClient = sns.NewFromConfig(cfg)

}

func GetSnsClient() *sns.Client{
	return snsClient
}
type SNSMessage struct{
	Email string
	Token string
}
func PublishMessage(email string,token string)error{
	err:=godotenv.Load(".env")
	if(err!=nil){
		log.Error().Err(err).Msg("Error loading Env for sns")
		return err
	}

	snsTopicArn=os.Getenv("Sns_Topic_Arn")
	emailData := SNSMessage{Email: email,
	Token: token}
	message, err := json.Marshal(emailData)
	if err != nil {
		log.Printf("Failed to marshal email data: %v", err)
	}
	// message:=email

	output,err:=snsClient.Publish(context.TODO(),&sns.PublishInput{
		Message: aws.String(string(message)),
		TopicArn: aws.String(snsTopicArn),
	})

	if(err!=nil){
		log.Error().Err(err).Msg(fmt.Sprintf("Error publising message %v",err))
	}

	log.Info().Msg(fmt.Sprintf("Message sent with ID: %s\n", *output.MessageId))

	return nil
}
