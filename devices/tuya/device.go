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

func DatabaseStatus(){

}

func scheduleSet (s* Schedule, device string) (error){
	key := device+"_schedule"
	bytes, err := json.Marshal(s)
	if err != nil{
		fmt.Println(err)
	}

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return err
	}
	defer c.Close()
	if _, err := c.Do("SET", key, string(bytes)); err != nil{
			fmt.Printf("Unable to set %s schedule\n", device)
			return err
	}

	return nil
}

func ScheduleGet (device string) (Schedule, error){
	s := Schedule{}
	key := device+"_schedule"

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
		return s, err
	}
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	json.Unmarshal([]byte(value), &s)
	return s, nil
}

func scheduleDel (device string) (error){
	key := device+"_schedule"

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

func scheduleUpdate (device string, status string) (error){
	s, err := ScheduleGet(device)
	if err != nil{
		fmt.Println(err)
		return err
	}

	s.Status = status

	if err := scheduleSet(&s,device); err != nil{
		fmt.Println(err)
		return err
	}

	return nil
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