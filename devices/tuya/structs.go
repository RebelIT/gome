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
	SundayOn 	string `json:"sunday_on"`
	SundayOff 	string `json:"sunday_off"`
	MondayOn	string `json:"monday_on"`
	MondayOff	string `json:"monday_off"`
	TuesdayOn	string `json:"tuesday_on"`
	TuesdayOff	string `json:"tuesday_off"`
	WednesdayOn string `json:"wednesday_on"`
	WednesdayOff	string `json:"wednesday_off"`
	ThursdayOn	string `json:"thursday_on"`
	ThursdayOff	string `json:"thursday_off"`
	FridayOn	string `json:"friday_on"`
	FridayOff	string `json:"friday_off"`
	SaturdayOn	string `json:"saturday_on"`
	SaturdayOff	string `json:"saturday_off"`
}

type ScheduleStatus struct{
	Enabled	bool `json:"enabled"`
}