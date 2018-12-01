package notify


type SlackMsg struct{
	Text		string `json:"text"`
	Username 	string `json:"username"`
	IconPath	string `json:"icon_path"`
}

type Secrets struct{
	SlackSecret		string `json:"slack_secret"`
}