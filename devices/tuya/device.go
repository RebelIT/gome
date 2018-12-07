package tuya

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"github.com/rebelit/gome/common"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func DeviceStatus (db string, ip string, id string, key string, name string) {
	fmt.Println("[DEBUG] Starting Device Status for "+name)
	data := Status{}

	args := []string{"get","--ip", ip,"--id", id, "--key", key}
	cmdOut, err := tryTuyaCli(string("tuya-cli"), args)
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

func scheduleSet (s* Schedules, device string) (error){
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

func ScheduleGet (device string) (Schedules, error){
	s := Schedules{}
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

func StatusGet (device string) (Status, error){
	s := Status{}
	key := device+"_status"

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
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

func detailsGet (device string) (Devices, error){
	d := Devices{}
	key := device

	c, err := dbConn()
	if err != nil{
		fmt.Println("Unable to connect to database")
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

func PowerControl(device string, value bool) error {
	d, err := detailsGet(device)
	if err != nil{
		return err
	}
	fmt.Printf("[INFO] issuing power control for %s\n", device)
	args := []string{"set","--id", d.Id, "--key", d.Key, "--set", strconv.FormatBool(value)}
	cmdOut, err := tryTuyaCli(string("tuya-cli"), args)
	if err != nil{
		return err
	} else {
		fmt.Printf("[DEBUG]: cmd return for %s : %s\n", device, cmdOut)
		fmtOut := strings.Replace(cmdOut, "\n", "", -1)
		if fmtOut == "Set succeeded."{
			return nil
		} else{
			return fmt.Errorf("error setting device status\n")
		}
	}
}

func tryTuyaCli(cmdName string, args []string) (string, error){
	maxRetry := 10
	retrySleep := time.Second * 1

	for i := 0; ;i++ {
		if i >= maxRetry{
			break
		}
		cmdOut, err := tuyaCli(cmdName, args)
		if err == nil{
			return cmdOut, err
		}
		fmt.Printf("[WARN] cmd %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %i reached for %s\n", maxRetry, cmdName)
}

func tuyaCli(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil{
		return "",err
	} else {
		fmtOut := strings.Replace(string(out), "\n", "", -1)
		if fmtOut == "Set succeeded." || fmtOut == "false" || fmtOut == "true" {
			return fmtOut, nil
		} else{
			return "", fmt.Errorf("error with tuya-cli\n")
		}
	}

}
func dbConn()(redis.Conn, error){
	var in Inputs

	deviceFile, err := ioutil.ReadFile(common.FILE)
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