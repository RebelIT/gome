package runner

import (
	"fmt"
	"github.com/rebelit/gome/devices"
	"github.com/rebelit/gome/devices/tuya"
	"github.com/rebelit/gome/notify"
	"log"
	"strconv"
	"strings"
	"time"
)

func GoGoScheduler() error {
	fmt.Println("[INFO] scheduler, starting")
	//var in Inputs

	for{
		devs, err := devices.LoadDevices()
		if err != nil{
			log.Printf("[WARN] scheduler, unable to load devices from file. skipping this round")
		}else {
			for _, d := range (devs.Devices) {
				switch d.Device {
				case "tuya":
					go doSchedule(d.Name)

				default:
					log.Printf("[WARN] scheduler, %s no device types match", d.Name)
				}
			}
		}
		time.Sleep(time.Second *60)
	}

	notify.SendSlackAlert("[ERROR] scheduler, routine broke out of loop")
	return nil
}

func doSchedule(device string) {
	_, iTime, day, _ := splitTime()
	schedule, err := devices.ScheduleGet(device)
	if err != nil {
		log.Printf("[WARN] scheduler, %s : %s\n", device, err)
	}

	devStatus, err := devices.StatusGet(device)
	if err != nil {
		log.Printf("[ERROR] scheduler, %s : %s\n", device, err)
	}

	scheduleOutCol := []Validator{}
	for _, s := range schedule.Schedules {
		if day == strings.ToLower(s.Day) {
			if s.Status == "enable" {
				scheduleOut := Validator{}

				onTime, _ := strconv.Atoi(s.On)   //time of day device is on
				offTime, _ := strconv.Atoi(s.Off) //time of day device is off

				//technically don't need `changeTo` anymore but I left it for now.
				doChange, changeTo, isInScheduleBlock := whatDoIDo(devStatus.Alive, iTime, onTime, offTime)

				scheduleOut.ChangeTo = changeTo
				scheduleOut.DoChange = doChange
				scheduleOut.InSchedule = isInScheduleBlock

				scheduleOutCol = append(scheduleOutCol, scheduleOut)
			}
		}
	}

	//if device is in any enabled schedule it must be on
	for _, s := range scheduleOutCol {
		if s.InSchedule && s.DoChange {
			if err := tuya.PowerControl(device, true); err != nil { //change it to true
				log.Printf("[ERROR] scheduler, %s failed to change powerstate: %s\n", device, err)
				notify.SendSlackAlert("[ERROR] scheduler failed to change powerstate for " + device)
			}
		}
	}

	//if device is not in  any enabled schedule it must be off
	if noSchedules(scheduleOutCol) {
		if devStatus.Alive {
			if err := tuya.PowerControl(device, false); err != nil { //change it to true
				log.Printf("[ERROR] scheduler, %s failed to change powerstate: %s\n", device, err)
				notify.SendSlackAlert("[ERROR] scheduler failed to change powerstate for " + device)
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

func noSchedules(v []Validator) bool {
	for _, s := range v {
		if s.InSchedule {
			return false
		}
	}
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