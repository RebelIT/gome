package common

import (
	"encoding/json"
	"fmt"
	"log"
)

func SendSlackAlert (message string){
	s, err := GetSecrets()
	if err != nil{
		log.Printf("[ERROR] slack alert: %s\n", err)
		return
	}

	content := SlackMsg{}
	content.Text = message
	content.Username = "gome"
	content.IconPath = ""

	body, _ := json.Marshal(content)
	resp, err := HttpPost( "https://hooks.slack.com/services/"+ s.SlackSecret, body,nil)
	if err != nil{
		log.Printf("[ERROR] slack alert: %s\n", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("[ERROR] slack alert: %s\n", fmt.Errorf("slack returned a non 200 response"))
		return
	} else{
		log.Printf("[INFO] slack alert sent: %s\n", message)
	}
}
