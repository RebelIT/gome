package devices

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices/rpi"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"io/ioutil"
	"log"
)


// *****************************************************************
// General device functions
func StatusGet (device string) (Status, error){
	s := Status{}

	values, err := DbHashGet(device+"_status")
	if err != nil{
		return s, err
	}
	redis.ScanStruct(values, &s)
	return s, nil
}

func DetailsGet (device string) (Devices, error){
	d := Devices{}

	values, err := DbHashGet(device)
	if err != nil{
		return d, err
	}
	redis.ScanStruct(values, &d)
	return d, nil
}

func LoadDevices()(Inputs, error){
	var in Inputs
	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		return in, err
	}
	json.Unmarshal(deviceFile, &in)

	return in, nil
}

func DoWhatAlexaSays(deviceType string, deviceName string, deviceAction string) error{
	action := false

	notify.MetricAws("alexa", "doAction", "nil",deviceName, deviceAction)

	switch deviceType{
	case "tuya":
		if deviceAction == "on"{
			action = true
		}
		if err := tuya.PowerControl(deviceName, action); err != nil{
			return err
		}
		return nil

	case "pi":
		if err := rpi.PiPost(deviceName, deviceAction); err != nil{
			return err
		}

	default:
		//no match
		return errors.New("no message in queue to parse")
	}

	return nil
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

func ScheduleGet (device string) (Schedules, error){
	s := Schedules{}

	value, err := DbGet(device+"_schedule")
	if err != nil{
		return s, err
	}
	json.Unmarshal([]byte(value), &s)

	if len(s.Schedules) <= 1 {
		return s, errors.New("invalid schedule struct")
	}

	return s, nil
}

func ScheduleDel (device string) (error){
	if err := DbDel(device+"_schedule"); err != nil{
		return err
	}
	return nil
}

func ScheduleUpdate (device string, status string) (error){
	s, err := ScheduleGet(device)
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