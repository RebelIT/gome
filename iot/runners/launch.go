package runners

import (
	"github.com/rebelit/gome/iot/runners/cron"
	"github.com/rebelit/gome/iot/runners/inventory"
	"time"
)

func Launch() {
	go inventory.GomeDevices()
	time.Sleep(time.Second * 2)

	go cron.GomeSchedules()
	time.Sleep(time.Second * 5)
}
