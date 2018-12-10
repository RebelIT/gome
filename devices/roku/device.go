package roku

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/notify"
	"io/ioutil"
	"log"
	"net/http"
)

func Connect(address string) (Roku, error) {
	return Roku{address: address, client: &http.Client{}}, nil
}

func (r Roku) Do(uriPart string) (http.Response, error) {
	url := "http://"+r.address+":8060"+uriPart

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return http.Response{}, err
	}
	resp, err := r.client.Do(req)
	return *resp, err
}

func (r Roku) Query(uriPart string) (http.Response, error) {
	url := "http://"+r.address+":8060"+uriPart

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.Response{},err
	}
	resp, err := r.client.Do(req)

	return *resp, err
}

func (r Roku) goHomeScreen() error{
	uri := "/keypress/home"
	resp, err := r.Do(uri)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func launchApp(deviceName string, app string) error {
	d, err := detailsGet(deviceName)
	if err != nil{
		return err
	}

	r, err := Connect(d.Addr)
	if err != nil {
		return err
	}

	id, err := getAppId(app)
	if err != nil{
	}

	uri := "/launch/"+id
	resp, err := r.Do(uri)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return err
	}

	return nil
}

func DeviceStatus(addr string, deviceName string){
	data := Status{}
	uriPart := "/"

	r, err := Connect(addr)
	if err != nil {
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	resp, err := r.Query(uriPart)
	if err != nil {
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}
	defer resp.Body.Close()

	notify.MetricHttpOut(deviceName, resp.StatusCode, "GET")

	if resp.StatusCode != 200 {
		data.Alive = false
	} else {
		data.Alive = true
	}
	data.Device = deviceName

	c, err := dbConn()
	if err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{deviceName+"_"+"status"}.AddFlat(data)); err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s : status done\n", deviceName)
	return
}

func detailsGet (device string) (Devices, error){
	d := Devices{}
	key := device

	c, err := dbConn()
	if err != nil{
		return Devices{}, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return Devices{}, err
	}

	redis.ScanStruct(values, &d)
	return d, nil
}

func StatusGet (device string) (Status, error){
	s := Status{}
	key := device+"_status"

	c, err := dbConn()
	if err != nil{
		return s, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return s, err
	}

	redis.ScanStruct(values, &s)

	return s, nil
}

func dbConn()(redis.Conn, error){
	var in Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(deviceFile, &in)

	db := in.Database
	conn, err := redis.Dial("tcp", db)
	if err != nil {
		return nil, err
	}
	return conn, nil
}