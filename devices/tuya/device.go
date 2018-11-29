package tuya

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

//Used for Device runner for running is alive inventory stored in redis
func DeviceStatus (db string, ip string, id string, key string, name string) {
	fmt.Println("[DEBUG] Starting Device Status for "+name)
	data := Status{}

	args := []string{"get","--ip", ip,"--id", id, "--key", key}
	cmdOut, err := command(string("tuya-cli"), args)
	if err != nil{
		fmt.Println("[ERROR] Error in tyua Cli, will Retry")
	}

	data.Device = name
	if strings.Replace(cmdOut, "\n", "", -1) == "true"{
		data.Alive = true
	} else {
		data.Alive = false
	}

	if err := cache.SetHash(db, redis.Args{name+"_"+"status"}.AddFlat(data)); err != nil {
		fmt.Println("[ERROR] Error in adding "+name+" to cache will retry")
		return
	}
	fmt.Println("[DEBUG] Done with Device Status for "+name)
	return
}

func scheduleGet (device string) (Schedule, error){
	s := Schedule{}
	key := device+"_schedule_details"

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return s, err
	}
	defer c.Close()

	values, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		fmt.Printf("unable to read %s from database\n", key)
		return s, err
	}

	redis.ScanStruct(values, &s)

	return s, nil
}

func scheduleSet (s* Schedule, device string) (error){
	key := device+"_schedule_details"
	dbData := redis.Args{key}.AddFlat(s)
	fmt.Printf("Data: %+v\n",dbData)
	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return err
	}
	defer c.Close()

	if _, err := c.Do("HMSET", dbData...); err != nil{
		fmt.Printf("Unable to set %s schedule\n", device)
		return err
	}

	return nil
}

func scheduleDel (device string) (error){
	key := device+"_schedule_details"

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return err
	}
	defer c.Close()

	if _, err := c.Do("DEL", key, "*"); err != nil{
		fmt.Printf("Unable to delete %s schedule\n", device)
		return err
	}
	return nil
}

func scheduleStatusSet (device string, status int) (error){
	key := device+"_schedule"

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return err
	}
	defer c.Close()

	c.Do("SET", key, status)
	return nil
}

func scheduleStatusGet(device string) (ScheduleStatus, error){
	key := device+"_schedule"
	s := ScheduleStatus{}

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return s, err
	}
	defer c.Close()

	value, err := redis.Bool(c.Do("GET", key))
	if err != nil {
		fmt.Printf("Unable to validate %s schedule\n", key)
		return s, err
	}
	s.Enabled = value
	return s, nil
}

func command(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil {
		log.Fatal(err)
		return "cmd error", err
	}

	return string(out), nil
}

func dbConn()(redis.Conn, error){
	var in Inputs

	deviceFile, err := ioutil.ReadFile(FILE)
	if err != nil {
		fmt.Println(err)
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