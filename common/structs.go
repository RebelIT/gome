package common

type Configs struct {
	redis 	string `json:"redis"`
	deviceFile 	string `json:"device_file"`
}

type Secrets struct{
	SlackSecret		string `json:"slack_secret"`
	AwsId 			string `json:"aws_id"`
	AwsSecret 		string `json:"aws_secret"`
	AwsToken		string `json:"aws_token"`
	AwsRegion 		string `json:"aws_region"`
	AWSQueueUrl 	string `json:"aws_queue_url"`
	RpiotUser		string `json:"rpiot_user"`
	RpiotToken		string `json:"rpiot_token"`
	Database 		string `json:"database"`
}

type SlackMsg struct{
	Text		string `json:"text"`
	Username 	string `json:"username"`
	IconPath	string `json:"icon_path"`
}