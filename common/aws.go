package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func AwsQueueClient() (sqs.SQS, string,  error){
	s, err := GetSecrets()
	if err != nil{
		return sqs.SQS{}, "", err
	}

	config, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.AwsRegion),
		Credentials: credentials.NewStaticCredentials(s.AwsAkid, s.AwsKey, s.AwsToken),
	})
	if err != nil{
		return sqs.SQS{}, "", err
	}

	return *sqs.New(config),s.AWSQueueUrl, nil
}