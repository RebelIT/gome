package tuya

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rebelit/gome/common"
	"github.com/rebelit/gome/devices"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func DeviceStatus (deviceName string, collectionDelayMin time.Duration) {
	log.Printf("[INFO] %s device collection delayed +%d sec\n",deviceName, collectionDelayMin)
	time.Sleep(time.Second * collectionDelayMin)

	doStatus := true
	alive := false

	d, err := devices.DetailsGet(deviceName+"_device")
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	cliArgs, err := generateGetCliArgs(d)
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
		doStatus = false
	}

	cmdOut, err := tryTuyaCli("tuya-cli", cliArgs)
	if err != nil{
		log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
	}

	if doStatus{
		if strings.Replace(cmdOut, "\n", "", -1) == "true"{
			alive = true
		}

		if err := devices.DbSet(deviceName+"_"+"status", []byte(strconv.FormatBool(alive))); err != nil{
			log.Printf("[ERROR] %s : device status, %s\n", deviceName, err)
			return
		}
	}

	log.Printf("[INFO] %s device status : done\n", deviceName)
	return
}

// device wrappers
func PowerControl(deviceName string, value bool) error {
	d, err := devices.DetailsGet(deviceName+"_device")
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
			if err := devices.UpdateStatus(deviceName, value); err != nil{
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

func generateSetCliArgs(deviceDetails devices.Devices, pwrState bool)(cliArg []string, err error){
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

func generateGetCliArgs(deviceDetails devices.Devices)(cliArg []string, err error){
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