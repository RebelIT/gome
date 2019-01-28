package rpi

type PiControl struct{
	Service   	string `json:"service"`
	Package		string `json:"package"`
	PinNumber 	string `json:"pin_number"`
	Action 		string `json:"action"`

}