package config

type Configuration struct {
	DbPath     string `json:"db_path"`
	StatsdAddr string `json:"statsd_addr"`
	HttpPort   string `json:"http_port"`
	SlackUrl   string `json:"slack_url"`
	RpiotUser  string `json:"rpiot_user"`
	RpiotToken string `json:"rpiot_token"`
	GomeUser   string `json:"gome_user"`
	GomePass   string `json:"gome_pass"`
}
