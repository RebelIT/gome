package cmd

import (
	"fmt"
	"github.com/rebelit/gome/util/notify"
	"log"
	"os/exec"
	"time"
)

func TryCommand(cmdName string, args []string, maxRetry int) (string, error) {
	retrySleep := time.Second * 2

	for i := 0; ; i++ {
		if i >= maxRetry {
			break
		}
		cmdOut, err := Command(cmdName, args)
		if err == nil {
			return cmdOut, err
		}
		notify.MetricCmd(cmdName, "retry")
		log.Printf("WARN: TryCommand, %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %d reached for %s\n", maxRetry, cmdName)
}

func Command(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil {
		notify.MetricCmd(cmdName, "failed")
		return string(out), err
	}
	notify.MetricCmd(cmdName, "success")
	return string(out), nil
}
