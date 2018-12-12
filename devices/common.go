package devices

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/common"
	"io/ioutil"
	"log"
)

func StatusGet (device string) (Status, error){
	s := Status{}
	key := device+"_status"

	c, err := DbConnect()
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

func DetailsGet (device string) (Devices, error){
	d := Devices{}
	key := device

	c, err := DbConnect()
	if err != nil{
		return Devices{}, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return Devices{}, err
	}
	redis.ScanStruct(values, &d)
	fmt.Printf("detailsGet db data %+v\n", d)
	return d, nil
}

func ScheduleSet (s* Schedules, device string) (error){
	key := device+"_schedule"
	bytes, err := json.Marshal(s)
	if err != nil{
		log.Println(err)
	}

	c, err := DbConnect()
	if err != nil{
		return err
	}
	defer c.Close()
	if _, err := c.Do("SET", key, string(bytes)); err != nil{
		return err
	}

	return nil
}

func ScheduleGet (device string) (Schedules, error){
	s := Schedules{}
	key := device+"_schedule"

	c, err := DbConnect()
	if err != nil{
		return s, err
	}
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	json.Unmarshal([]byte(value), &s)
	return s, nil
}

func ScheduleDel (device string) (error){
	key := device+"_schedule"

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

func ScheduleUpdate (device string, status string) (error){
	s, err := ScheduleGet(device)
	if err != nil{
		fmt.Println(err)
		return err
	}

	s.Status = status

	if err := ScheduleSet(&s,device); err != nil{
		fmt.Println(err)
		return err
	}

	return nil
}

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

func LoadDevices()(Inputs, error){
	var in Inputs
	deviceFile, err := ioutil.ReadFile(common.FILE)
	if err != nil {
		fmt.Println(err)
		return in, err
	}
	json.Unmarshal(deviceFile, &in)

	return in, nil
}