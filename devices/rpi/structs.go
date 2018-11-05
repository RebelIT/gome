package rpi

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

type DeviceAction struct{
	DeviceName	string `json:"device_name"`
	Action		string `json:"action"`
}

type RtnOptRoot struct{
	ActionUri		string
}

type RtnOptAction struct{
	Actions		[]Action
}

type Action struct{
	Uri		string	`json:"endpoint"`
	Method	string	`json:"method"`
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