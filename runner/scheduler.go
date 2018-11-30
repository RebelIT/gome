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
	fmt.Println("[DEBUG] Starting scheduler")
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
				fmt.Printf("[WARN] No device type match for %s, dont schedule anything", d.Name)
			}
		}
		time.Sleep(time.Second *60)
	}

	return nil
}

func doSchedule(device string) error {
	sTime, iTime, day, now := splitTime()

	s, err := tuya.ScheduleGet(device)
	if err != nil{
		fmt.Println(err)
	}

	on := 0
	off := 0
	if s.Status == "enable"{
		//TODO: get status before proceeding to compare it in case/if. 
		tuya.DatabaseStatus()

		switch day {
		case "sunday":
			on, _= strconv.Atoi(s.Days.Sunday.On)
			off, _= strconv.Atoi(s.Days.Sunday.Off)

			if on <= iTime{
				if off >= iTime{
					//turon dev on
				}
			}


		case "monday":

		case "tuesday":

		case "wednesday":

		case "thursday":

		case "friday":

		case "saturday":

		default:
			fmt.Println("[ERROR] we've entered the matrix where days of the week no longer matter")
			panic("This can't be right...")
		}
	}

	return nil
}
func splitTime()(strTime string, intTime int, weekday string, now time.Time){
	Now := time.Now()
	NowMinute := Now.Minute()
	NowHour := Now.Hour()
	NowDay := Now.Weekday()

	sTime := strconv.Itoa(NowHour) + strconv.Itoa(NowMinute)
	iTime, _ := strconv.Atoi(sTime)
	day := strings.ToLower(NowDay.String())

	return sTime, iTime, day, Now
}
