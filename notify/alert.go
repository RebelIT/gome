package notify

import (
	"fmt"
	"github.com/rebelit/gome/common"
	"log"
)


func SendSlackAlert (message string){
	s, err := getSecrets()
	if err != nil{
		log.Printf("[ERROR] slack alert: %s\n", err)
		return
	}

	content := SlackMsg{}
	content.Text = message
	content.Username = "gome"

	respCode, err := common.PostWebReq(content, "https://hooks.slack.com/services/"+ s.SlackSecret)
	if err != nil{
		log.Printf("[ERROR] slack alert: %s\n", err)
		return
	}
	if respCode != 200 {
		log.Printf("[ERROR] slack alert: %s\n", fmt.Errorf("slack returned a non 200 response"))
		return
	} else{
		log.Printf("[INFO] slack alert sent: %s\n", message)
	}
	MetricHttpOut("slack", respCode, "POST")
}
