package devices

type Inputs struct{
	Database	string `json:"database"`
	Devices		[]Devices
}

type Devices struct{
	Device			string `json:"device"`
	Type 			string `json:"type"`
	Name			string `json:"name"`
	Addr 			string `json:"address"`
	NameFriendly	string `json:"name_friendly"`
	NetPort			string `json:"port"`
	Id				string `json:"id"`
	Key 			string `json:"key"`
	Dps 			string `json:"dps"`
}

type Status struct{
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

type DeviceAction struct{
	Action		bool `json:"action"`
}

type PiControl struct{
	Service   	string `json:"service"`
	Package		string `json:"package"`
	PinNumber 	string `json:"pin_number"`
	Action 		string `json:"action"`
}

type PiResponse struct {
	Namespace 	string `json:"namespace"`
	Message		string `json:"message"`
}