package listener

import "net/http"

//core handlers to update the device inventory.  devices.json used on gome startup to load last state into redis
//

func getDevices(w http.ResponseWriter,r *http.Request){
	//GET - return all known devices

}

func addDevice(w http.ResponseWriter,r *http.Request){
	//POST - add row to json and redis

}

func delDevice(w http.ResponseWriter,r *http.Request){
	//DELETE - row from json and redis

}

func updateDevice(w http.ResponseWriter,r *http.Request){
	//PUT - update json and redis

}
