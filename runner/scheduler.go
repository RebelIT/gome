package runner

import (
	"encoding/json"
	"fmt"
	"github.com/rebelit/gome/devices/tuya"
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
		time.Sleep(time.Second *30)
	}

	return nil
}

func doSchedule(device string) error {
	fmt.Printf("[DEBUG] Schedule for %s\n", device)
	_, iTime, day, _ := splitTime()

	s, err := tuya.ScheduleGet(device)
	if err != nil{
		fmt.Printf("[WARN] could not get schedule or schedule does not exist yet\n")
		fmt.Printf("[WARN] %s\n", err)
	}
	fmt.Printf("[DEBUG] Schedule Status for %s %s\n", device, s.Status)
	if s.Status == "enable"{
		status, err := tuya.StatusGet(device)
		if err != nil{
			fmt.Printf("[ERROR] getting device status from database: %s\n", err)
		}

		devState := status.Alive
		fmt.Printf("[DEBUG] Device Alive Status for %s %v\n", device, devState)
		//fmt.Printf("[DEBUG] Day %s\n",day)

		switch day {
		case "sunday":
			onTime, _:= strconv.Atoi(s.Days.Sunday.On) //time of day device is on
			offTime, _:= strconv.Atoi(s.Days.Sunday.Off) //time of day device is off

			fmt.Printf("[DEBUG] Device on/off for %s %v %v\n",device, onTime, offTime)
			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange {
				fmt.Printf("[DEBUG] Changing Status %v : change to %v\n",doChange, powerState)
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "monday":
			onTime, _:= strconv.Atoi(s.Days.Monday.On)
			offTime, _:= strconv.Atoi(s.Days.Monday.Off)

			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "tuesday":
			onTime, _:= strconv.Atoi(s.Days.Tuesday.On)
			offTime, _:= strconv.Atoi(s.Days.Tuesday.Off)

			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "wednesday":
			onTime, _:= strconv.Atoi(s.Days.Wednesday.On)
			offTime, _:= strconv.Atoi(s.Days.Wednesday.Off)

			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "thursday":
			onTime, _:= strconv.Atoi(s.Days.Thursday.On)
			offTime, _:= strconv.Atoi(s.Days.Thursday.Off)

			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "friday":
			onTime, _:= strconv.Atoi(s.Days.Friday.On)
			offTime, _:= strconv.Atoi(s.Days.Friday.Off)

			fmt.Printf("[DEBUG] Device on/off for %s %v %v\n",device, onTime, offTime)
			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange {
				fmt.Printf("[DEBUG] Changing Status %v : change to %v\n",doChange, powerState)
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s", err)
					return err
				}
			}
			return nil

		case "saturday":
			onTime, _:= strconv.Atoi(s.Days.Saturday.On)
			offTime, _:= strconv.Atoi(s.Days.Saturday.Off)

			doChange, powerState := whatDoIDo(devState, iTime, onTime, offTime)

			if doChange{
				if err := tuya.PowerControl(device, powerState);err != nil{
					fmt.Printf("[ERROR] failed to change powerstate: %s\n", err)
					return err
				}
			}
			return nil

		default:
			fmt.Println("[ERROR] we've entered the matrix where days of the week no longer exist")
			panic("This can't be right...")
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
	}

	sTime = strconv.Itoa(NowHour) + strconv.Itoa(NowMinute)
	iTime, _ := strconv.Atoi(sTime)
	day := strings.ToLower(NowDay.String())

	//fmt.Printf("[DEBUG] sTime: %v\n",sTime)
	//fmt.Printf("[DEBUG] iTime: %v\n",iTime)
	//fmt.Printf("[DEBUG] day: %v\n",day)
	//fmt.Printf("[DEBUG] Now: %v\n",Now)
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