package rpi

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

type DeviceAction struct{
	Action		string `json:"action"`
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