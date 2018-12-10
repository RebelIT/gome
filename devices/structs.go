package devices

type Inputs struct{
	Database	string `json:"database"`
	Devices		[]Devices
}

type Devices struct{
	Device		string `json:"device"`
	Name		string `json:"name"`
	Addr 		string `json:"address"`
	NetPort		string `json:"port"`
	Id			string `json:"id"`
	Key 		string `json:"key"`
}

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

type DeviceAction struct{
	Action		bool `json:"action"`
}

type Schedules struct {
	Status 		string `json:"status"`
	Schedules		[]Schedule `json:"schedules"`
}

type Schedule struct{
	Day		string `json:"day"`
	Status 	string `json:"status"`
	Desc	string `json:"desc"`
	On		string `json:"on"`
	Off 	string `json:"off"`
}