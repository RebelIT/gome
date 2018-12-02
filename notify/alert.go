package notify

import (
	"fmt"
	"github.com/rebelit/gome/common"
)

const FILE  = "/etc/gome/devices.json"

func SendSlackAlert (message string){
	s, err := getSecrets()
	if err != nil{
		fmt.Printf("[ERROR] slack alert: %s\n", err)
	}

	content := SlackMsg{}
	content.Text = message
	content.Username = "gome"

	respCode, err := common.PostWebReq(content, "https://hooks.slack.com/services/"+ s.SlackSecret)
	if err != nil{
		fmt.Printf("[ERROR] slack alert: %s\n", err)
	}
	if respCode != 200 {
		fmt.Printf("[ERROR] slack alert: %s\n", fmt.Errorf("slack returned a non 200 response"))
	}
}
