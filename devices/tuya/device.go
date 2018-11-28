package tuya

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/rebelit/gome/cache"
	"log"
	"os/exec"
	"strings"
)

//Used for Device runner for running is alive inventory stored in redis
func DeviceStatus(db string, id string, key string, name string) {
	fmt.Println("[DEBUG] Starting Device Status for "+name)
	data := Status{}

	args := []string{"get","--id", id, "--key", key}
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

func command(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil {
		log.Fatal(err)
		return "cmd error", err
	}

	return string(out), nil
}