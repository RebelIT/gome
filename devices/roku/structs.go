package roku

import "net/http"

type Roku struct {
	address string
	client  *http.Client
}

type Devices struct{
	Device		string `json:"device"`
	Name		string `json:"name"`
	Addr 		string `json:"address"`
	NetPort		string `json:"port"`
	Id			string `json:"id"`
	Key 		string `json:"key"`
}

type Inputs struct{
	Database	string `json:"database"`
	Devices		[]Devices
}

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

//XML :facepalm
type App struct {
	App		string `xml:",chardata"`
	Id		string  `xml:"id,attr"`
}

type Apps struct {
	Apps	[]App		`xml:"app"`
}

type DeviceAction struct{
	DeviceName	string `json:"device_name"`
	Action		string `json:"action"`
	ActionItem	string `json:"action_item"`
}