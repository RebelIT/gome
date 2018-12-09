package notify

import (
	"github.com/rebelit/gome/common"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
	"strconv"
)


func MetricDeviceStatus(deviceName string, deviceType string, alive bool){
	//emits a new counter for every external web request
	tags := statsd.Tags("device", deviceName, "device_type", deviceType)
	measurement := "gome_device"
	value := ""

	if alive{
		value = "1"
	} else{
		value = "0"
	}

	sendUnique(measurement, tags, value)
}

func MetricHttpIn(uri string, reponseCode int, method string){
	//emits a new counter for every incoming http web request
	tags := statsd.Tags("uri", uri, "status_code", strconv.Itoa(reponseCode))
	measurement := "gome_http"

	sendCounter(measurement, tags)
}

func MetricHttpOut(uri string, reponseCode int, method string){
	//emits a new counter for every external web request
	tags := statsd.Tags("uri", uri, "response_code", strconv.Itoa(reponseCode))
	measurement := "gome_http_out"

	sendCounter(measurement, tags)
}

func sendCounter(measurement string, tags statsd.Option){
	addrOpt := statsd.Address(common.STATSD)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	c, err := statsd.New(addrOpt,fmtOpt,tags)
	if err != nil {
		log.Print(err)
	}
	defer c.Close()

	c.Increment(measurement)
}

func sendUnique(measurement string, tags statsd.Option, value string){
	addrOpt := statsd.Address(common.STATSD)
	fmtOpt := statsd.TagsFormat(statsd.InfluxDB)
	c, err := statsd.New(addrOpt,fmtOpt,tags)
	if err != nil {
		log.Print(err)
	}
	defer c.Close()

	c.Unique(measurement, value)
}