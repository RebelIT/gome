package runner

type Devices struct{
	Device			string `json:"device"`
	Name			string `json:"name"`
	NameFriendly	string `json:"name_friendly"`
	Addr 			string `json:"address"`
	NetPort			string `json:"port"`
	Id				string `json:"id"`
	Key 			string `json:"key"`
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
	InSchedule 	bool
}