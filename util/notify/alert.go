package notify

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome/util/config"
	"github.com/rebelit/gome/util/http"
	"log"
)

func Slack(message string) {
	content := SlackMsg{}
	content.Text = message
	content.Username = "gome"
	content.IconPath = ""

	body, _ := json.Marshal(content)
	resp, err := http.HttpPost("https://hooks.slack.com/services/"+config.App.SlackUrl, body, nil)
	if err != nil {
		log.Printf("WARN: Slack, %s\n", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Printf("WARN: Slack, %s\n", fmt.Errorf("slack returned a non 200 response"))
		return
	}

	log.Printf("[INFO] Slack, sent: %s\n", message)
	return
}
