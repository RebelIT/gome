package runner

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

//type Validator struct{
//	States []State
//}

type Validator struct{
	DoChange 	bool
	ChangeTo	bool
	InSchedule 	bool
}