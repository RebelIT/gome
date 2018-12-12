package tuya

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/notify"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func DeviceStatus (addr string, id string, key string, deviceName string) {
	data := devices.Status{}

	args := []string{"get","--ip", addr,"--id", id, "--key", key}
	cmdOut, err := tryTuyaCli(string("tuya-cli"), args)
	if err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		notify.MetricCmd("tuya-cli", "failed")
	}
	notify.MetricCmd("tuya-cli", "success")

	if strings.Replace(cmdOut, "\n", "", -1) == "true"{
		data.Alive = true
	} else {
		data.Alive = false
	}
	data.Device = deviceName

	c, err := devices.DbConnect()
	if err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{deviceName+"_"+"status"}.AddFlat(data)...); err != nil{
		log.Printf("[ERROR] %s : status, %s\n", deviceName, err)
		return
	}

	log.Printf("[DEBUG] %s : status done\n", deviceName)
	return
}

func PowerControl(device string, value bool) error {
	d, err := devices.DetailsGet(device)
	if err != nil{
		return err
	}
	log.Printf("[INFO] issuing power control for %s\n", device)
	args := []string{"set","--ip", d.Addr, "--id", d.Id, "--key", d.Key, "--set", strconv.FormatBool(value)}
	cmdOut, err := tryTuyaCli("tuya-cli", args)
	if err != nil{
		return err
	} else {
		log.Printf("[DEBUG]: cmd return for %s : %s\n", device, cmdOut)
		fmtOut := strings.Replace(cmdOut, "\n", "", -1)
		if fmtOut == "Set succeeded."{
			notify.SendSlackAlert("Tuya PowerControl initiated for "+device+" to "+strconv.FormatBool(value)+"")
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
		notify.MetricCmd("tuya-cli", "retry")
		log.Printf("[WARN] cmd %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %i reached for %s\n", maxRetry, cmdName)
}

func tuyaCli(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil{
		notify.MetricCmd("tuya-cli", "failed")
		return "",err
	} else {
		fmtOut := strings.Replace(string(out), "\n", "", -1)
		if fmtOut == "Set succeeded." || fmtOut == "false" || fmtOut == "true" {
			notify.MetricCmd("tuya-cli", "success")
			return fmtOut, nil
		} else{
			notify.MetricCmd("tuya-cli", "failed")
			return "", fmt.Errorf("error with tuya-cli\n")
		}
	}

}