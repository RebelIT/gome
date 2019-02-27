package common

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func TryCommand(cmdName string, args []string, maxRetry int) (string, error){
	retrySleep := time.Second * 2

	for i := 0; ;i++ {
		if i >= maxRetry{
			break
		}
		cmdOut, err := Command(cmdName, args)
		if err == nil{
			return cmdOut, err
		}
		MetricCmd(cmdName, "retry")
		log.Printf("[WARN] cmd %s failed, retrying\n", cmdName)
		time.Sleep(retrySleep)
	}

	return "", fmt.Errorf("max retries %d reached for %s\n", maxRetry, cmdName)
}

func Command(cmdName string, args []string) (string, error) {
	out, err := exec.Command(cmdName, args...).Output()
	if err != nil {
		MetricCmd(cmdName, "failed")
		return string(out), err
	}
	MetricCmd(cmdName, "success")
	return string(out), nil
}