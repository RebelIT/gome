package runner

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"log"
	"strings"
	"time"
)

func GoGoSQS() error {
	log.Println("[INFO] aws sqs, starting")
	notify.SendSlackAlert("AWS SQS runner is starting")

	s, err := common.GetSecrets()
	if err != nil {
		return err
	}

	config, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.AwsRegion),
		Credentials: credentials.NewStaticCredentials(s.AwsId, s.AwsSecret, s.AwsToken),
	})
	if err != nil {
		return err
	}

	c := sqs.New(config)

	for {
		message, receipt, err := getMessage(c, s.AWSQueueUrl)
		if err != nil {
			log.Printf("[WARN] aws sqs, %s", err)
			if receipt != nil{
				if err := deleteMessage(c,s.AWSQueueUrl,receipt); err != nil{
					log.Printf("[WARN] aws sqs, %s", err)
				}
			}
		} else {
			m := strings.Split(message, ",")
			deviceType := m[0]
			deviceName := m[1]
			deviceAction := m[2]

			if err := deleteMessage(c,s.AWSQueueUrl,receipt); err != nil{
				log.Printf("[WARN] aws sqs, %s", err)
			}

			if err := doAction(deviceType,deviceName,deviceAction); err != nil{
				log.Printf("[ERROR], aws sqs, %s", err)
			}
		}
		time.Sleep(time.Second *2)
	}
	notify.SendSlackAlert("AWS SQS runner broke out of the loop. Get it back in there")
	return nil
}

func doAction(deviceType string, deviceName string, deviceAction string) error{
	action := false

	notify.MetricAws("sqs", "doAction", "nil",deviceName, deviceAction)

	switch deviceType{
	case "tuya":
		if deviceAction == "on"{
			action = true
		}
		if err := tuya.PowerControl(deviceName, action); err != nil{
			return err
		}
		return nil

	default:
		//no match
		return errors.New("no message in queue to parse")
	}

	return nil
}

func getMessage(c *sqs.SQS, queueUrl string)(string, *string, error){
	message := ""

	param := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl), // Required
		AttributeNames: []*string{
			aws.String("QueueAttributeName"), // Required
			// More values...
		},
		MaxNumberOfMessages: aws.Int64(5),
		MessageAttributeNames: []*string{
			aws.String("MessageAttributeName"), // Required
			// More values...
		},
		VisibilityTimeout: aws.Int64(10),
		WaitTimeSeconds:   aws.Int64(0),
	}

	result, err := c.ReceiveMessage(param)
	if err != nil{
		notify.MetricAws("sqs", "get", "failure","nil", "nil")
		return "", nil, err
	}
	notify.MetricAws("sqs", "get", "success","nil", "nil")

	if len(result.Messages) == 0 {
		return "", nil, errors.New("no messages in queue")
	} else {
		message = awsutil.StringValue(result.Messages[0].Body)
		if message == ""{
			return "", result.Messages[0].ReceiptHandle, errors.New("message was blank when it should have not been...")
		}
	}
	return message, result.Messages[0].ReceiptHandle, nil
}

func deleteMessage(c *sqs.SQS,queueUrl string, receipt *string) error{
	param := &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: receipt,
	}
	_, err := c.DeleteMessage(param)
	if err != nil {
		notify.MetricAws("sqs", "delete", "failure","nil", "nil")
		return err
	}
	notify.MetricAws("sqs", "delete", "success","nil", "nil")
	return nil
}
