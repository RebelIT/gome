package listener

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

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
