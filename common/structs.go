package common

type Configs struct {
	redis 	string `json:"redis"`
	deviceFile 	string `json:"device_file"`
}

type Secrets struct{
	SlackSecret		string `json:"slack_secret"`
	AwsAkid 		string `json:"aws_akid"`
	AwsKey 			string `json:"aws_key"`
	AwsToken		string `json:"aws_token"`
	AwsRegion 		string `json:"aws_region"`
	AWSQueueUrl 	string `json:"aws_queue_url"`
}