package tuya

type Status struct{
	Device	string `json:"device"`
	Alive	bool `json:"alive"`
	Url		string `json:"url"`
}

type DeviceAction struct{
	Action		bool `json:"action"`
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

type Schedule struct {
	Status		string `json:"status"`
	Days		ScheduleDay `json:"days"`
}

type ScheduleDay struct{
	Sunday 		ScheduleOnOff `json:"sunday"`
	Monday		ScheduleOnOff `json:"monday"`
	Tuesday		ScheduleOnOff `json:"tuesday"`
	Wednesday 	ScheduleOnOff `json:"wednesday"`
	Thursday	ScheduleOnOff `json:"thursday"`
	Friday		ScheduleOnOff `json:"friday"`
	Saturday	ScheduleOnOff `json:"saturday"`
}

type ScheduleOnOff struct{
	On	string `json:"on"`
	Off string `json:"off"`
}