package devices

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/database"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func TuyaDeviceStatus (deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	doStatus := true
	alive := false
	powerState := false

	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	cliArgs, err := generateGetCliArgs(d)
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	if doStatus{
		cmdOut, err := tryTuyaCli("tuya-cli", cliArgs)
		if err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		}else{
			alive = true
		}

		if strings.Replace(cmdOut, "\n", "", -1) == "true"{
			powerState = true
		}

		//status = is on the network and accessible
		if err := database.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}

		//state = the current switch power state on/true | off/false
		if err := database.DbSet(deviceName+"_power"+"_state", []byte(strconv.FormatBool(powerState))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

// device wrappers
func TuyaPowerControl(deviceName string, value bool) error {
	d, err := DetailsGet(deviceName+"_device")
	if err != nil{
		return err
	}

	cliArgs, err := generateSetCliArgs(d, value)
	if err != nil{
		return err
	}

	log.Printf("[INFO] issuing power control for %s\n", deviceName)
	cmdOut, err := tryTuyaCli("tuya-cli", cliArgs)
	if err != nil{
		return err
	} else {
		fmtOut := strings.Replace(cmdOut, "\n", "", -1)
		if fmtOut == "Set succeeded."{
			if err := UpdateStatus(deviceName, value); err != nil{
				log.Printf("[ERROR] Update Device Status, %s : %s", deviceName, err)
			}
			common.SendSlackAlert("Tuya PowerControl initiated for "+d.Name+"("+d.NameFriendly+") to "+strconv.FormatBool(value)+"")
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
		common.MetricCmd("tuya-cli", "retry")
		log.Printf("[WARN] cmd %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %i reached for %s\n", maxRetry, cmdName)
}

func tuyaCli(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil{
		common.MetricCmd("tuya-cli", "failed")
		return "",err
	} else {
		fmtOut := strings.Replace(string(out), "\n", "", -1)
		if fmtOut == "Set succeeded." || fmtOut == "false" || fmtOut == "true" {
			common.MetricCmd("tuya-cli", "success")
			return fmtOut, nil
		} else{
			common.MetricCmd("tuya-cli", "failed")
			return "", fmt.Errorf("error with tuya-cli\n")
		}
	}

}

func generateSetCliArgs(deviceDetails Devices, pwrState bool)(cliArg []string, err error){
	args := []string{}
	fmt.Printf("%+v\n",deviceDetails)
	switch deviceDetails.Type{
		case "outlet":
			args = []string{"set", "--set", strconv.FormatBool(pwrState), "--ip", deviceDetails.Addr,
							"--id", deviceDetails.Id, "--key", deviceDetails.Key}

		case "switch":
			args = []string{"set", "--set", strconv.FormatBool(pwrState), "--ip", deviceDetails.Addr,
				"--id", deviceDetails.Id, "--key", deviceDetails.Key, "--dps", deviceDetails.Dps}

		default:
			return args, errors.New("no device type "+deviceDetails.Type+" found in cli args switch")
	}
	return args, nil
}

func generateGetCliArgs(deviceDetails Devices)(cliArg []string, err error){
	args := []string{}

	switch deviceDetails.Type{
	case "outlet":
		args = []string{"get", "--ip", deviceDetails.Addr, "--id", deviceDetails.Id, "--key", deviceDetails.Key}

	case "switch":
		args = []string{"get", "--ip", deviceDetails.Addr, "--id", deviceDetails.Id, "--key", deviceDetails.Key,
						"--dps", deviceDetails.Dps}

	default:
		return args, errors.New("no device type "+deviceDetails.Type+" found in cli args switch")
	}
	return args, nil
}