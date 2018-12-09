package notify

import (
	"github.com/rebelit/gome/common"
	"gopkg.in/alexcesaro/statsd.v2"
	"strconv"
)


func MetricDeviceStatus(deviceName string, deviceType string, alive bool){
	//emits a new generic metric for device status
	tags := statsd.Tags("device", deviceName, "device_type", deviceType)
	measurement := "gome_device"
	value := ""

	if alive{
		value = "on"
	} else{
		value = "off"
	}

	sendUnique(measurement, tags, value)
}

func MetricHttpIn(uri string, reponseCode int, method string){
	//emits a new counter for every incoming http web request
	tags := statsd.Tags("uri", uri, "status_code", strconv.Itoa(reponseCode),"method", method)
	measurement := "gome_http"

	sendCounter(measurement, tags)
}

func MetricHttpOut(destination string, reponseCode int, method string){
	//emits a new counter for every external web request
	tags := statsd.Tags("destination", destination, "response_code", strconv.Itoa(reponseCode), "method", method)
	measurement := "gome_http_out"

	sendCounter(measurement, tags)
}

func MetricCmd(cmd string, response string){
	//emits a new counter for every shell out function
	tags := statsd.Tags("cmd", cmd, "cmd_status", response)
	measurement := "gome_cmd"

	sendCounter(measurement, tags)
}

func sendCounter(measurement string, tags statsd.Option){
	if common.STATSD == ""{
		//metrics disabled don't do anything
		return
	}
	addrOpt := statsd.Address(common.STATSD)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	c, err := statsd.New(addrOpt,fmtOpt,tags)
	if err != nil {
		//log.Print(err)
	}
	defer c.Close()

	c.Increment(measurement)
}

func sendUnique(measurement string, tags statsd.Option, value string){
	if common.STATSD == ""{
		//metrics disabled don't do anything
		return
	}
	addrOpt := statsd.Address(common.STATSD)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	c, err := statsd.New(addrOpt,fmtOpt,tags)
	if err != nil {
		//log.Print(err)
	}
	defer c.Close()
	c.Unique(measurement, value)
}