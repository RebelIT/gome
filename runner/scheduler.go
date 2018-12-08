package runner

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome/common"
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
		deviceFile, err := ioutil.ReadFile(common.FILE)
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

func doSchedule(device string) {
	_, iTime, day, _ := splitTime()
	schedule, err := tuya.ScheduleGet(device)
	if err != nil {
		fmt.Printf("[WARN] could not get schedule or schedule does not exist yet\n")
		fmt.Printf("[WARN] %s\n", err)
	}

	devStatus, err := tuya.StatusGet(device)
	if err != nil {
		fmt.Printf("[ERROR] getting device status from database: %s\n", err)
	}

	scheduleOutCol := []Validator{}
	for _, s := range schedule.Schedules {
		if day == strings.ToLower(s.Day) {
			if s.Status == "enable" {
				scheduleOut := Validator{}

				onTime, _ := strconv.Atoi(s.On)   //time of day device is on
				offTime, _ := strconv.Atoi(s.Off) //time of day device is off

				//technically don't need `doChange` and `changeTo` anymore but I left it to debug later for mismatch issues.
				doChange, changeTo, isInScheduleBlock := whatDoIDo(devStatus.Alive, iTime, onTime, offTime)

				scheduleOut.ChangeTo = changeTo
				scheduleOut.DoChange = doChange
				scheduleOut.InSchedule = isInScheduleBlock

				scheduleOutCol = append(scheduleOutCol, scheduleOut)
			}
		}
		fmt.Printf("[VERBOSE] %s full validate array:  %+v\n", device, scheduleOutCol)
		//if device is in any enabled schedule it must be on
		for _, s := range scheduleOutCol {
			fmt.Printf("[VERBOSE] %s inSchedule %v\n", device, s.InSchedule)
			if s.InSchedule {
				fmt.Printf("[VERBOSE] %s inSchedule %v validated changing it\n", device, s.InSchedule)
				if err := tuya.PowerControl(device, true); err != nil { //change it to true
					fmt.Printf("[ERROR] failed to change powerstate: %s\n", err)
					notify.SendSlackAlert("Scheduler [ERROR] failed to change powerstate for " + device)
				}
			}
		}

		//if device is not in  any enabled schedule it must be off
		if !isInAllSchedules(scheduleOutCol) {
			fmt.Printf("[VERBOSE] %s evaluates not in any schedule\n", device)
			if err := tuya.PowerControl(device, false); err != nil { //change it to true
				fmt.Printf("[ERROR] failed to change powerstate: %s\n", err)
				notify.SendSlackAlert("Scheduler [ERROR] failed to change powerstate for " + device)
			}
		}
	}
}

func splitTime()(strTime string, intTime int, weekday string, now time.Time){
	Now := time.Now()
	NowMinute := Now.Minute()
	NowHour := Now.Hour()
	NowDay := Now.Weekday()

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

func isInAllSchedules(v []Validator) bool {
	fmt.Printf("[VERBOSE] parsing isInAllSchedules\n")
	if len(v) > 1 {
		a := v[0].InSchedule
		fmt.Printf("[VERBOSE] parsing isInAllSchedules 1st item %v\n", v[0].InSchedule)
		for _, s := range v {
			fmt.Printf("[VERBOSE] parsing isInAllSchedules each item %v\n", s.InSchedule)
			if a != s.InSchedule {
				fmt.Printf("[VERBOSE] parsing isInAllSchedules each item evaluated false match\n")
				return false
			}
		}
	}
	fmt.Printf("[VERBOSE] parsing isInAllSchedules each item evaluated true match\n")
	return true
}

func whatDoIDo(devOn bool, currentHour int, devOnTime int, devOffTime int) (changeState bool, changeTo bool, inScheduleBlock bool){
	reverseCheck := false
	changeState = false
	changeTo = false
	inScheduleBlock = false

	if devOffTime <= devOnTime {
		//spans a day PM to AM on schedule
		reverseCheck = true
	}

	if !reverseCheck{
		//does not span PM to AM
		inScheduleBlock = inBetween(currentHour, devOnTime, devOffTime)
	} else {
		//spans a day PM to AM reverse check the schedule
		inScheduleBlock = inBetweenReverse(currentHour, devOffTime, devOnTime)
	}

	if devOn{
		if inScheduleBlock{
			//leave it be change state:false
			changeState = false
			return changeState, changeTo, inScheduleBlock
		} else {
			//change state:true. change the power control to false
			changeState = true
			changeTo = false
			return changeState, changeTo, inScheduleBlock
		}
	} else {
		if inScheduleBlock{
			//change state:true. change the power control to true
			changeState = true
			changeTo = true
			return changeState, changeTo, inScheduleBlock
		}else {
			//leave it be change state:false
			changeState = false
			return changeState, changeTo, inScheduleBlock
		}
	}
}