package runner

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices/tuya"
	"log"
	"time"
)

func GoGoSQS() error {
	log.Println("[INFO] aws sqs, starting")
	log.Println("[DEBUG] aws sqs, loading new session")

	s, err := common.GetSecrets()
	if err != nil {
		return err
	}

	config, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.AwsRegion),
		Credentials: credentials.NewStaticCredentials(s.AwsAkid, s.AwsKey, s.AwsToken),
	})
	if err != nil {
		return err
	}

	c := sqs.New(config)

	for {
		result, err := c.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &s.AWSQueueUrl,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(20), // 20 seconds
			WaitTimeSeconds:     aws.Int64(0),
		})
		if err != nil {
			log.Printf("[ERROR] aws sqs, message : %s\n", err)
			return err
		}

		if len(result.Messages) == 0 {
			fmt.Println("Received no messages")
			return err
		} else {
			for _, m := range result.Messages{
				if err := parseMessage(*m); err != nil{
					break
				}

				resultDelete, err := c.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &s.AWSQueueUrl,
					ReceiptHandle: m.ReceiptHandle,
				})
				if err != nil {
					fmt.Println("Delete Error", err)
					return err
				}
				log.Printf("[INFO] aws msg delete, %s\n", resultDelete.String())
			}

		}
	}
	time.Sleep(time.Second * 2)

	return nil
}

func parseMessage(message sqs.Message) error{
	body := message.String()
	switch body{
	case "gome_test":
		if err := tuya.PowerControl("treeFamily", false); err != nil{
			return err
		}
		return nil

	default:
		//no match
		return errors.New("no message in queue to parse")
	}

	return nil
}