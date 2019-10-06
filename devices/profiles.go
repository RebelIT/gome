package devices

//Universal Smart Things Profiles
type Devices struct {
	Device []Device `json:"device"`
}

type Device struct {
	Make        string   `json:"make"`
	Model       string   `json:"model"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Region      string   `json:"region"`
	Location    string   `json:"location"`
	Metadata    Metadata `json:"meta"`
	Actions     []string `json:"actions"`
	State       State    `json:"state"`
}

type Metadata struct {
	NetAddr  string `json:"net_addr"`
	Port     string `json:"port"`
	Id       string `json:"id"`
	Key      string `json:"key"`
	Dps      string `json:"dps"`
	Pin      string `json:"pin"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type State struct {
	Alive  bool `json:"alive"`
	Status bool `json:"status"`
}

type Inputs struct {
	Devices []DevicesOld
}

type DevicesOld struct {
	Device       string `json:"device"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Addr         string `json:"address"`
	NameFriendly string `json:"name_friendly"`
	NetPort      string `json:"port"`
	Id           string `json:"id"`
	Key          string `json:"key"`
	Dps          string `json:"dps"`
}

type Status struct {
	Alive bool   `json:"alive"`
	Url   string `json:"url"`
}

type DeviceAction struct {
	Action bool `json:"action"`
}

type PiControl struct {
	Service   string `json:"service"`
	Package   string `json:"package"`
	PinNumber string `json:"pin_number"`
	Action    string `json:"action"`
}

type PiResponse struct {
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
}
