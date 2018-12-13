package runner

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rebelit/gome/common"
	"log"
)

func GoGoSQS() error {
	log.Println("[INFO] aws sqs, starting")
	log.Println("[DEBUG] aws sqs, loading new session")

	c, qUrl, err := common.AwsQueueClient()
	if err != nil{
		fmt.Printf("[ERROR] aws sqs, unable to create client: %s\n", err)
	}

	//start for
	result, err := c.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20),  // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return err
	}

	resultDelete, err := c.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qUrl,
		ReceiptHandle: result.Messages[0].ReceiptHandle,
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		return err
	}

	fmt.Println("Message Deleted", resultDelete)

	//end for
	//sleep
	return nil
	//end function



	//for {
	//	log.Println("Checking Queue")
	//
	//}
	//time.Sleep(time.Second *2)

}

