package devices

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"io/ioutil"
	"log"
	"net/http"
)

// *****************************************************************
// General device functions
func StatusGet (device string) (status string, error error){  //Gets the device status from redis
	value, err := DbGet(device+"_status")
	if err != nil{
		return "", err
	}

	return value, nil
}

func UpdateStatus(deviceName string, status bool) error{  //Update the device status in redis
	statusData := Status{}
	statusData.Alive = status

	if err := DbHashSet(deviceName+"_"+"status", statusData); err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		return err
	}

	return nil
}

func DetailsGet (device string) (Devices, error){  //Gets the device details from redis
	d := Devices{}

	values, err := DbHashGet(device)
	if err != nil{
		return d, err
	}
	redis.ScanStruct(values, &d)
	return d, nil
}

func LoadDevices() error{  //Load Devices into redis from devices.json file
	log.Printf("[INFO] device loader, starting")
	i, err := ReadDeviceFile()
	if err != nil{
		return err
	}

	if len(i.Devices) == 0{
		log.Printf("[WARN] device loader, no devices to load, skipping")
		return nil
	}

	for _, d := range i.Devices {
		log.Printf("[INFO] device loader, loading %s under '%s_device'", d.Name, d.Name)
		if err := DbHashSet(d.Name+"_device",d); err != nil{
			return err
		}
	}
	log.Println("[INFO] device loader, all done")
	return nil
}

func ReadDeviceFile()(Inputs, error){  //Read the devices.json
	var in Inputs
	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		return in, err
	}
	if err := json.Unmarshal(deviceFile, &in); err != nil{
		return in, err
	}

	return in, nil
}

func GetAllDevicesFromDb() (devices []string, err error){  //Gets a full inventory of devices from redis
	keySearch := "*_device"
	keys, err := DbGetKeys(keySearch)
	if err != nil{
		return nil, err
	}
	return keys, nil
}


// *****************************************************************
// Scheduler functions
func ScheduleSet (s* Schedules, device string) (error){
	data, err := json.Marshal(s)
	if err != nil{
		log.Println(err)
	}

	if err := DbSet(device+"_schedule", data); err != nil{
		return err
	}
	return nil
}

func ScheduleGet (devName string) (hasSchedule bool, schedules Schedules, error error){
	s := Schedules{}

	value, err := DbGet(devName+"_schedule")
	if err != nil{
		return false, s, err
	}
	if value == ""{
		return false, s, nil
	}

	json.Unmarshal([]byte(value), &s)

	if len(s.Schedules) <= 1 {
		return false, s, errors.New("invalid schedule struct")
	}

	return true, s, nil
}

func ScheduleDel (device string) (error){
	if err := DbDel(device+"_schedule"); err != nil{
		return err
	}
	return nil
}

func ScheduleUpdate (device string, status string) (error){
	_, s, err := ScheduleGet(device)
	if err != nil{
		return err
	}

	s.Status = status

	if err := ScheduleSet(&s,device); err != nil{
		return err
	}

	return nil
}


// *****************************************************************
// Redis functions
func DbConnect()(redis.Conn, error){
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

func DbHashSet(key string, data interface{} ) error{
	//equivalent to a redis HMSET
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()
	if _, err := c.Do("HMSET", redis.Args{key}.AddFlat(data)...); err != nil{
		return err
	}

	return nil
}

func DbHashGet(key string)(values []interface{}, err error){
	c, err := DbConnect()
	if err != nil{
		return nil, err
	}
	defer c.Close()

	resp, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func DbGetKeys(key string)(keys []string, err error){
	c, err := DbConnect()
	if err != nil{
		return nil, err
	}
	defer c.Close()

	resp, err := redis.Strings(c.Do("KEYS", key))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DbSet(key string, value []byte) error{
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()
	if _, err := c.Do("SET", key, string(value)); err != nil{
		return err
	}
	return nil
}

func DbGet(key string) (values string, err error){
	c, err := DbConnect()
	if err != nil{
		return "", err
	}
	defer c.Close()
	value, err := redis.String(c.Do("GET", key))
	if err != nil{
		if err.Error() == "redigo: nil returned"{
			//This is fine, redis connection OK, just no data returned
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func DbDel(key string) error{
	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()

	if _, err := c.Do("DEL", key, "*"); err != nil{
		return err
	}

	return nil
}


// *****************************************************************
// Http Response helper functions
func ReturnOk(w http.ResponseWriter, r *http.Request, response interface{}){
	code := http.StatusOK
	common.MetricHttpIn(r.RequestURI, code, r.Method)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] %s : %s\n", r.URL.Path, err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	return
}

func ReturnBad(w http.ResponseWriter, r *http.Request){
	code := http.StatusBadRequest
	common.MetricHttpIn(r.RequestURI, code, r.Method)
	w.WriteHeader(code)
	return
}

func ReturnInternalError(w http.ResponseWriter, r *http.Request){
	code := http.StatusInternalServerError
	common.MetricHttpIn(r.RequestURI, code, r.Method)
	w.WriteHeader(code)
	return
}
