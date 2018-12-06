package runner

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GoGoScheduler() error {
	fmt.Println("[INFO] Starting scheduler")
	var in Inputs

	for {
		deviceFile, err := ioutil.ReadFile(FILE)
		if err != nil {
			fmt.Println(err)
			return err
		}
		json.Unmarshal(deviceFile, &in)

		for _, d := range (in.Devices) {
			switch d.Device {
			case "tuya":

				go doSchedule(d.Name)

			default:
				fmt.Printf("[WARN] No device type match for %s, dont schedule anything\n", d.Name)
			}
		}
		time.Sleep(time.Second *60)
	}

	notify.SendSlackAlert("Scheduler broke out of loop")
	return nil
}

func doSchedule(device string) error {
	_, iTime, day, _ := splitTime()
	schedule, err := tuya.ScheduleGet(device)
	if err != nil{
		fmt.Printf("[WARN] could not get schedule or schedule does not exist yet\n")
		fmt.Printf("[WARN] %s\n", err)
	}

	devStatus, err := tuya.StatusGet(device)
	if err != nil{
		fmt.Printf("[ERROR] getting device status from database: %s\n", err)
	}

	for _, s := range schedule.Schedules{
		if day == strings.ToLower(s.Day) && s.Status == "enable"{
			onTime, _:= strconv.Atoi(s.On)   //time of day device is on
			offTime, _:= strconv.Atoi(s.Off) //time of day device is off

			doChange, powerState := whatDoIDo(devStatus.Alive, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s\n", err)
					notify.SendSlackAlert("Scheduler [ERROR] failed to change powerstate for "+ device)
					return err
				}
				notify.SendSlackAlert("Scheduler "+device+" changed from "+strconv.FormatBool(devStatus.Alive)+" to "+strconv.FormatBool(powerState))
			}
			return nil

		}
	}

	return nil
}

func splitTime()(strTime string, intTime int, weekday string, now time.Time){
	Now := time.Now()
	NowMinute := Now.Minute()
	NowHour := Now.Hour()
	NowDay := now.Weekday()

	sTime := ""
	singleMinute := inBetween(NowMinute, 0,9)
	if singleMinute{
		sTime = strconv.Itoa(NowHour) + "0"+ strconv.Itoa(NowMinute)
	} else{
		sTime = strconv.Itoa(NowHour) + strconv.Itoa(NowMinute)
	}

	iTime, _ := strconv.Atoi(sTime)
	day := strings.ToLower(NowDay.String())

	return sTime, iTime, day, Now
}

func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

func inBetweenReverse(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return false
	} else {
		return true
	}
}

func whatDoIDo(devOn bool, currentHour int, devOnTime int, devOffTime int) (changeState bool, changeTo bool){
	reverseCheck := false
	ok := false
	changeState = false
	changeTo = false

	if devOffTime <= devOnTime {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck{
		//does not span PM to AM
		ok = inBetween(currentHour, devOnTime, devOffTime)
	} else {
		//spans a day PM to AM reverse check the schedule
		ok = inBetweenReverse(currentHour, devOffTime, devOnTime)
	}

	if devOn{
		if ok{
			//leave it be change state:false
			changeState = false
			return changeState, changeTo
		} else {
			//change state:true. change the power control to false
			changeState = true
			changeTo = false
			return changeState, changeTo
		}
	} else {
		if ok{
			//change state:true. change the power control to true
			changeState = true
			changeTo = true
			return changeState, changeTo
		}else {
			//leave it be change state:false
			changeState = false
			return changeState, changeTo
		}
	}
}